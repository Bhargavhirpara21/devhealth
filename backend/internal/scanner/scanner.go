package scanner

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/go-github/v62/github"

	"github.com/BhargavHirpara/devhealth/internal/models"
	"github.com/BhargavHirpara/devhealth/internal/scoring"
)

// Scanner scans GitHub repositories for health checks.
type Scanner struct {
	client *github.Client
}

// New creates a new Scanner with an authenticated GitHub client.
func New(client *github.Client) *Scanner {
	return &Scanner{client: client}
}

// ScanOwner scans all non-archived, non-fork repos for a given owner.
func (s *Scanner) ScanOwner(ctx context.Context, owner, ownerType string) ([]models.RepoHealth, error) {
	repos, err := s.listRepos(ctx, owner, ownerType)
	if err != nil {
		return nil, fmt.Errorf("listing repos: %w", err)
	}

	var results []models.RepoHealth
	for _, repo := range repos {
		if repo.GetArchived() || repo.GetFork() {
			continue
		}

		rh, err := s.ScanRepo(ctx, owner, repo.GetName(), repo)
		if err != nil {
			log.Printf("warning: failed to scan %s/%s: %v", owner, repo.GetName(), err)
			continue
		}
		results = append(results, *rh)
	}

	return results, nil
}

// ScanRepo runs all health checks on a single repository.
func (s *Scanner) ScanRepo(ctx context.Context, owner, repoName string, repo *github.Repository) (*models.RepoHealth, error) {
	if repo == nil {
		var err error
		repo, _, err = s.client.Repositories.Get(ctx, owner, repoName)
		if err != nil {
			return nil, fmt.Errorf("getting repo: %w", err)
		}
	}

	// Use canonical owner/repo from GitHub API to avoid case-sensitivity issues
	owner = repo.GetOwner().GetLogin()
	repoName = repo.GetName()

	defaultBranch := repo.GetDefaultBranch()
	if defaultBranch == "" {
		defaultBranch = "main"
	}

	checks := []models.CheckResult{
		s.checkBranchProtection(ctx, owner, repoName, defaultBranch),
		s.checkSecretScanning(repo),
		s.checkDependabot(ctx, owner, repoName),
		s.checkCICD(ctx, owner, repoName),
		s.checkReadme(ctx, owner, repoName),
		s.checkLicense(ctx, owner, repoName),
		s.checkCodeowners(ctx, owner, repoName),
		s.checkVulnerabilities(ctx, owner, repoName),
		s.checkStaleRepo(repo),
	}

	score := scoring.Calculate(checks)

	rh := &models.RepoHealth{
		Owner:         owner,
		Repo:          repoName,
		FullName:      fmt.Sprintf("%s/%s", owner, repoName),
		URL:           repo.GetHTMLURL(),
		Score:         score,
		Checks:        checks,
		ScannedAt:     time.Now().UTC(),
		DefaultBranch: defaultBranch,
		IsArchived:    repo.GetArchived(),
		IsFork:        repo.GetFork(),
	}

	if !repo.GetPushedAt().IsZero() {
		rh.LastCommitAt = repo.GetPushedAt().Time
	}

	return rh, nil
}

func (s *Scanner) listRepos(ctx context.Context, owner, ownerType string) ([]*github.Repository, error) {
	var allRepos []*github.Repository

	if ownerType == "org" {
		opts := &github.RepositoryListByOrgOptions{
			Type:        "all",
			ListOptions: github.ListOptions{PerPage: 100},
		}
		for {
			repos, resp, err := s.client.Repositories.ListByOrg(ctx, owner, opts)
			if err != nil {
				return nil, err
			}
			allRepos = append(allRepos, repos...)
			if resp.NextPage == 0 {
				break
			}
			opts.Page = resp.NextPage
		}
	} else {
		opts := &github.RepositoryListByUserOptions{
			Type:        "owner",
			ListOptions: github.ListOptions{PerPage: 100},
		}
		for {
			repos, resp, err := s.client.Repositories.ListByUser(ctx, owner, opts)
			if err != nil {
				return nil, err
			}
			allRepos = append(allRepos, repos...)
			if resp.NextPage == 0 {
				break
			}
			opts.Page = resp.NextPage
		}
	}

	return allRepos, nil
}

func (s *Scanner) checkBranchProtection(ctx context.Context, owner, repo, branch string) models.CheckResult {
	result := models.CheckResult{
		Name:           models.CheckBranchProtection,
		Severity:       "critical",
		Recommendation: "Enable branch protection rules on the default branch to prevent direct pushes and require pull request reviews.",
	}

	protection, _, err := s.client.Repositories.GetBranchProtection(ctx, owner, repo, branch)
	if err != nil {
		result.Passed = false
		result.Details = "Branch protection is not configured on the default branch."
		return result
	}

	result.Passed = true
	details := "Branch protection is enabled"
	if protection.RequiredPullRequestReviews != nil {
		details += fmt.Sprintf("; requires %d review(s)", protection.RequiredPullRequestReviews.RequiredApprovingReviewCount)
	}
	if protection.RequiredStatusChecks != nil {
		details += "; has required status checks"
	}
	result.Details = details
	return result
}

func (s *Scanner) checkSecretScanning(repo *github.Repository) models.CheckResult {
	result := models.CheckResult{
		Name:           models.CheckSecretScanning,
		Severity:       "critical",
		Recommendation: "Enable secret scanning in repository security settings to detect accidentally committed credentials.",
	}

	securityFeatures := repo.GetSecurityAndAnalysis()
	if securityFeatures != nil && securityFeatures.SecretScanning != nil {
		result.Passed = securityFeatures.SecretScanning.GetStatus() == "enabled"
		if result.Passed {
			result.Details = "Secret scanning is enabled."
		} else {
			result.Details = "Secret scanning is disabled."
		}
	} else {
		result.Passed = false
		result.Details = "Secret scanning status could not be determined."
	}

	return result
}

func (s *Scanner) checkDependabot(ctx context.Context, owner, repo string) models.CheckResult {
	result := models.CheckResult{
		Name:           models.CheckDependabot,
		Severity:       "high",
		Recommendation: "Enable Dependabot alerts to receive notifications about vulnerable dependencies.",
	}

	alerts, _, err := s.client.Dependabot.ListRepoAlerts(ctx, owner, repo, nil)
	if err != nil {
		// If we get a 404 or error, Dependabot may not be enabled.
		result.Passed = false
		result.Details = "Dependabot alerts are not enabled or accessible."
		return result
	}

	result.Passed = true
	openCount := 0
	for _, a := range alerts {
		if a.GetState() == "open" {
			openCount++
		}
	}
	result.Details = fmt.Sprintf("Dependabot is enabled; %d open alert(s).", openCount)
	return result
}

func (s *Scanner) checkCICD(ctx context.Context, owner, repo string) models.CheckResult {
	result := models.CheckResult{
		Name:           models.CheckCICD,
		Severity:       "high",
		Recommendation: "Add a GitHub Actions workflow (.github/workflows/) to automate testing and deployment.",
	}

	_, dirContent, _, err := s.client.Repositories.GetContents(ctx, owner, repo, ".github/workflows", nil)
	if err != nil {
		result.Passed = false
		result.Details = "No CI/CD pipeline found (.github/workflows/ does not exist)."
		return result
	}

	yamlCount := 0
	for _, f := range dirContent {
		name := f.GetName()
		if len(name) > 4 && (name[len(name)-4:] == ".yml" || name[len(name)-5:] == ".yaml") {
			yamlCount++
		}
	}

	result.Passed = yamlCount > 0
	if result.Passed {
		result.Details = fmt.Sprintf("Found %d workflow file(s) in .github/workflows/.", yamlCount)
	} else {
		result.Details = ".github/workflows/ exists but contains no workflow files."
	}
	return result
}

func (s *Scanner) checkReadme(ctx context.Context, owner, repo string) models.CheckResult {
	result := models.CheckResult{
		Name:           models.CheckReadme,
		Severity:       "medium",
		Recommendation: "Add a README.md to document the project's purpose, setup, and usage.",
	}

	fileContent, _, _, err := s.client.Repositories.GetContents(ctx, owner, repo, "README.md", nil)
	if err != nil {
		result.Passed = false
		result.Details = "No README.md found in the repository root."
		return result
	}

	result.Passed = true
	result.Details = fmt.Sprintf("README.md exists (%d bytes).", fileContent.GetSize())
	return result
}

func (s *Scanner) checkLicense(ctx context.Context, owner, repo string) models.CheckResult {
	result := models.CheckResult{
		Name:           models.CheckLicense,
		Severity:       "low",
		Recommendation: "Add a LICENSE file to clarify how the code can be used.",
	}

	license, _, err := s.client.Repositories.License(ctx, owner, repo)
	if err != nil {
		result.Passed = false
		result.Details = "No license file detected."
		return result
	}

	result.Passed = true
	if license.License != nil {
		result.Details = fmt.Sprintf("License: %s.", license.License.GetName())
	} else {
		result.Details = "License file exists."
	}
	return result
}

func (s *Scanner) checkCodeowners(ctx context.Context, owner, repo string) models.CheckResult {
	result := models.CheckResult{
		Name:           models.CheckCodeowners,
		Severity:       "low",
		Recommendation: "Add a CODEOWNERS file to define ownership and streamline code reviews.",
	}

	// CODEOWNERS can be in root, .github/, or docs/
	paths := []string{"CODEOWNERS", ".github/CODEOWNERS", "docs/CODEOWNERS"}
	for _, path := range paths {
		_, _, _, err := s.client.Repositories.GetContents(ctx, owner, repo, path, nil)
		if err == nil {
			result.Passed = true
			result.Details = fmt.Sprintf("CODEOWNERS file found at %s.", path)
			return result
		}
	}

	result.Passed = false
	result.Details = "No CODEOWNERS file found."
	return result
}

func (s *Scanner) checkVulnerabilities(ctx context.Context, owner, repo string) models.CheckResult {
	result := models.CheckResult{
		Name:           models.CheckVulnerabilities,
		Severity:       "high",
		Recommendation: "Review and resolve open vulnerability alerts to reduce security risk.",
	}

	alerts, _, err := s.client.Dependabot.ListRepoAlerts(ctx, owner, repo, &github.ListAlertsOptions{
		State: github.String("open"),
	})
	if err != nil {
		// Can't determine — give benefit of doubt
		result.Passed = true
		result.Details = "Could not check vulnerability alerts (may lack permissions)."
		return result
	}

	openCount := len(alerts)
	result.Passed = openCount == 0
	if result.Passed {
		result.Details = "No open vulnerability alerts."
	} else {
		result.Details = fmt.Sprintf("%d open vulnerability alert(s).", openCount)
	}
	return result
}

func (s *Scanner) checkStaleRepo(repo *github.Repository) models.CheckResult {
	result := models.CheckResult{
		Name:           models.CheckStaleRepo,
		Severity:       "low",
		Recommendation: "Consider archiving or updating this repository if it is no longer actively maintained.",
	}

	if repo.GetPushedAt().IsZero() {
		result.Passed = false
		result.Details = "No push activity recorded."
		return result
	}

	daysSinceLastPush := time.Since(repo.GetPushedAt().Time).Hours() / 24
	result.Passed = daysSinceLastPush < 180 // 6 months threshold

	if result.Passed {
		result.Details = fmt.Sprintf("Last push was %.0f day(s) ago.", daysSinceLastPush)
	} else {
		result.Details = fmt.Sprintf("Repository appears stale — last push was %.0f day(s) ago.", daysSinceLastPush)
	}
	return result
}
