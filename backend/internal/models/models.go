package models

import "time"

// CheckName identifies a specific health check.
type CheckName string

const (
	CheckBranchProtection CheckName = "branch_protection"
	CheckSecretScanning   CheckName = "secret_scanning"
	CheckDependabot       CheckName = "dependabot_alerts"
	CheckCICD             CheckName = "cicd_pipeline"
	CheckReadme           CheckName = "readme"
	CheckLicense          CheckName = "license"
	CheckCodeowners       CheckName = "codeowners"
	CheckVulnerabilities  CheckName = "vulnerabilities"
	CheckStaleRepo        CheckName = "stale_repo"
)

// CheckResult holds the outcome of a single health check.
type CheckResult struct {
	Name           CheckName `json:"name"`
	Passed         bool      `json:"passed"`
	Details        string    `json:"details"`
	Severity       string    `json:"severity"` // "critical", "high", "medium", "low"
	Recommendation string    `json:"recommendation"`
}

// RepoHealth holds the full health report for a single repository.
type RepoHealth struct {
	ID          int64         `json:"id"`
	Owner       string        `json:"owner"`
	Repo        string        `json:"repo"`
	FullName    string        `json:"full_name"`
	URL         string        `json:"url"`
	Score       int           `json:"score"`
	Checks      []CheckResult `json:"checks"`
	ScannedAt   time.Time     `json:"scanned_at"`
	DefaultBranch string     `json:"default_branch"`
	LastCommitAt  time.Time  `json:"last_commit_at"`
	IsArchived    bool       `json:"is_archived"`
	IsFork        bool       `json:"is_fork"`
}

// ScanRequest is the payload for triggering a scan.
type ScanRequest struct {
	Owner string `json:"owner"`
	Type  string `json:"type"` // "org" or "user"
}

// ScanResponse is returned after triggering a scan.
type ScanResponse struct {
	Message    string `json:"message"`
	ReposFound int    `json:"repos_found"`
	ScanID     string `json:"scan_id"`
}

// Summary holds aggregate stats for a scanned owner.
type Summary struct {
	Owner          string         `json:"owner"`
	TotalRepos     int            `json:"total_repos"`
	AverageScore   float64        `json:"average_score"`
	ScoreDistribution ScoreBuckets `json:"score_distribution"`
	CheckPassRates map[CheckName]float64 `json:"check_pass_rates"`
	LastScanAt     time.Time      `json:"last_scan_at"`
}

// ScoreBuckets groups repos by score ranges.
type ScoreBuckets struct {
	Critical int `json:"critical"` // 0-39
	Warning  int `json:"warning"`  // 40-69
	Good     int `json:"good"`     // 70-89
	Excellent int `json:"excellent"` // 90-100
}

// ErrorResponse is a standard API error.
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}
