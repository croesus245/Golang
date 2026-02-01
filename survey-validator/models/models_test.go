package models

import (
	"testing"
)

func TestSurveyPoint_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		point    SurveyPoint
		expected bool
	}{
		{
			name:     "Valid point",
			point:    SurveyPoint{PointID: "P1", Easting: 100, Northing: 200},
			expected: true,
		},
		{
			name:     "Empty point ID",
			point:    SurveyPoint{PointID: "", Easting: 100, Northing: 200},
			expected: false,
		},
		{
			name:     "Zero easting",
			point:    SurveyPoint{PointID: "P1", Easting: 0, Northing: 200},
			expected: false,
		},
		{
			name:     "Zero northing",
			point:    SurveyPoint{PointID: "P1", Easting: 100, Northing: 0},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.point.IsValid()
			if result != tt.expected {
				t.Errorf("IsValid() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestSurveyPoint_HasHeight(t *testing.T) {
	height := 100.5

	withHeight := SurveyPoint{PointID: "P1", Easting: 100, Northing: 200, Height: &height}
	if !withHeight.HasHeight() {
		t.Error("Expected HasHeight() to be true")
	}

	withoutHeight := SurveyPoint{PointID: "P1", Easting: 100, Northing: 200}
	if withoutHeight.HasHeight() {
		t.Error("Expected HasHeight() to be false")
	}
}

func TestValidationReport_AddIssue(t *testing.T) {
	report := NewValidationReport("TEST")

	if report.Status != StatusPass {
		t.Errorf("Initial status should be PASS, got %s", report.Status)
	}

	// Add warning
	report.AddIssue(ValidationIssue{
		Severity:    SeverityWarning,
		Description: "Test warning",
	})

	if report.Status != StatusWarning {
		t.Errorf("Status should be WARNING after adding warning, got %s", report.Status)
	}

	// Add error
	report.AddIssue(ValidationIssue{
		Severity:    SeverityError,
		Description: "Test error",
	})

	if report.Status != StatusFail {
		t.Errorf("Status should be FAIL after adding error, got %s", report.Status)
	}

	if len(report.Issues) != 2 {
		t.Errorf("Expected 2 issues, got %d", len(report.Issues))
	}
}

func TestValidationReport_CalculateConfidenceScore(t *testing.T) {
	report := NewValidationReport("TEST")
	report.CalculateConfidenceScore()

	if report.ConfidenceScore != 100 {
		t.Errorf("Empty report should have score 100, got %f", report.ConfidenceScore)
	}

	// Add issues
	report.AddIssue(ValidationIssue{Severity: SeverityError})   // -15
	report.AddIssue(ValidationIssue{Severity: SeverityWarning}) // -5
	report.AddIssue(ValidationIssue{Severity: SeverityInfo})    // -1
	report.CalculateConfidenceScore()

	expected := 100.0 - 15.0 - 5.0 - 1.0
	if report.ConfidenceScore != expected {
		t.Errorf("Expected score %f, got %f", expected, report.ConfidenceScore)
	}
}
