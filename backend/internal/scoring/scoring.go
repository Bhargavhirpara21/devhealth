package scoring

import "github.com/BhargavHirpara/devhealth/internal/models"

// Weight defines a check's weight toward the total score.
type Weight struct {
	Name   models.CheckName
	Points int
}

// DefaultWeights defines how much each check contributes to the total score.
// Total possible points = 100.
var DefaultWeights = []Weight{
	{models.CheckBranchProtection, 20}, // Critical: security
	{models.CheckSecretScanning, 15},   // Critical: security
	{models.CheckDependabot, 15},       // High: vulnerability management
	{models.CheckCICD, 15},             // High: quality assurance
	{models.CheckVulnerabilities, 10},  // High: active vulnerabilities
	{models.CheckReadme, 10},           // Medium: documentation
	{models.CheckLicense, 5},           // Low: compliance
	{models.CheckCodeowners, 5},        // Low: governance
	{models.CheckStaleRepo, 5},         // Low: maintenance
}

// Calculate calculates the health score (0–100) from check results.
func Calculate(checks []models.CheckResult) int {
	if len(checks) == 0 {
		return 0
	}

	passed := make(map[models.CheckName]bool)
	for _, c := range checks {
		passed[c.Name] = c.Passed
	}

	score := 0
	for _, w := range DefaultWeights {
		if passed[w.Name] {
			score += w.Points
		}
	}

	return score
}
