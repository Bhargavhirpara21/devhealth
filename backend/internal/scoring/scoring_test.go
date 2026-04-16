package scoring

import (
	"testing"

	"github.com/BhargavHirpara/devhealth/internal/models"
)

func TestCalculate_AllPassed(t *testing.T) {
	checks := []models.CheckResult{
		{Name: models.CheckBranchProtection, Passed: true},
		{Name: models.CheckSecretScanning, Passed: true},
		{Name: models.CheckDependabot, Passed: true},
		{Name: models.CheckCICD, Passed: true},
		{Name: models.CheckVulnerabilities, Passed: true},
		{Name: models.CheckReadme, Passed: true},
		{Name: models.CheckLicense, Passed: true},
		{Name: models.CheckCodeowners, Passed: true},
		{Name: models.CheckStaleRepo, Passed: true},
	}

	score := Calculate(checks)
	if score != 100 {
		t.Errorf("expected 100, got %d", score)
	}
}

func TestCalculate_NonePassed(t *testing.T) {
	checks := []models.CheckResult{
		{Name: models.CheckBranchProtection, Passed: false},
		{Name: models.CheckSecretScanning, Passed: false},
		{Name: models.CheckDependabot, Passed: false},
		{Name: models.CheckCICD, Passed: false},
		{Name: models.CheckVulnerabilities, Passed: false},
		{Name: models.CheckReadme, Passed: false},
		{Name: models.CheckLicense, Passed: false},
		{Name: models.CheckCodeowners, Passed: false},
		{Name: models.CheckStaleRepo, Passed: false},
	}

	score := Calculate(checks)
	if score != 0 {
		t.Errorf("expected 0, got %d", score)
	}
}

func TestCalculate_OnlyCriticalPassed(t *testing.T) {
	checks := []models.CheckResult{
		{Name: models.CheckBranchProtection, Passed: true},
		{Name: models.CheckSecretScanning, Passed: true},
		{Name: models.CheckDependabot, Passed: false},
		{Name: models.CheckCICD, Passed: false},
		{Name: models.CheckVulnerabilities, Passed: false},
		{Name: models.CheckReadme, Passed: false},
		{Name: models.CheckLicense, Passed: false},
		{Name: models.CheckCodeowners, Passed: false},
		{Name: models.CheckStaleRepo, Passed: false},
	}

	score := Calculate(checks)
	expected := 35 // 20 + 15
	if score != expected {
		t.Errorf("expected %d, got %d", expected, score)
	}
}

func TestCalculate_Empty(t *testing.T) {
	score := Calculate([]models.CheckResult{})
	if score != 0 {
		t.Errorf("expected 0, got %d", score)
	}
}

func TestCalculate_MixedResults(t *testing.T) {
	checks := []models.CheckResult{
		{Name: models.CheckBranchProtection, Passed: true},  // 20
		{Name: models.CheckSecretScanning, Passed: false},    // 0
		{Name: models.CheckDependabot, Passed: true},         // 15
		{Name: models.CheckCICD, Passed: true},               // 15
		{Name: models.CheckVulnerabilities, Passed: false},   // 0
		{Name: models.CheckReadme, Passed: true},             // 10
		{Name: models.CheckLicense, Passed: true},            // 5
		{Name: models.CheckCodeowners, Passed: false},        // 0
		{Name: models.CheckStaleRepo, Passed: true},          // 5
	}

	score := Calculate(checks)
	expected := 70 // 20 + 15 + 15 + 10 + 5 + 5
	if score != expected {
		t.Errorf("expected %d, got %d", expected, score)
	}
}
