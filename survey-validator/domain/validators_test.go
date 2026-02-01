package domain

import (
	"testing"

	"github.com/survey-validator/models"
)

func TestValidateInput_EmptyPoints(t *testing.T) {
	data := &models.SurveyData{
		ProjectID: "TEST",
		Points:    []models.SurveyPoint{},
	}

	issues := ValidateInput(data)

	if len(issues) != 1 {
		t.Errorf("Expected 1 issue, got %d", len(issues))
	}

	if issues[0].Severity != models.SeverityError {
		t.Errorf("Expected error severity, got %s", issues[0].Severity)
	}
}

func TestValidateInput_EmptyPointID(t *testing.T) {
	data := &models.SurveyData{
		ProjectID: "TEST",
		Points: []models.SurveyPoint{
			{PointID: "", Easting: 100, Northing: 100},
		},
	}

	issues := ValidateInput(data)

	hasEmptyIDIssue := false
	for _, issue := range issues {
		if issue.Description == "Point found with empty Point ID" {
			hasEmptyIDIssue = true
			break
		}
	}

	if !hasEmptyIDIssue {
		t.Error("Expected issue for empty point ID")
	}
}

func TestDetectDuplicates(t *testing.T) {
	data := &models.SurveyData{
		ProjectID: "TEST",
		Points: []models.SurveyPoint{
			{PointID: "P1", Easting: 100, Northing: 100},
			{PointID: "P2", Easting: 100.0005, Northing: 100.0005}, // Near duplicate
			{PointID: "P3", Easting: 200, Northing: 200},           // Not duplicate
		},
	}

	issues := DetectDuplicates(data)

	if len(issues) == 0 {
		t.Error("Expected at least one duplicate issue")
	}

	// Should detect P1 and P2 as duplicates or near-duplicates
	found := false
	for _, issue := range issues {
		if len(issue.PointIDs) == 2 {
			if (issue.PointIDs[0] == "P1" && issue.PointIDs[1] == "P2") ||
				(issue.PointIDs[0] == "P2" && issue.PointIDs[1] == "P1") {
				found = true
				break
			}
		}
	}

	if !found {
		t.Error("Expected P1 and P2 to be flagged as duplicates")
	}
}

func TestDetectOutliers(t *testing.T) {
	// Create a cluster of points with one clear outlier
	// The outlier needs to be far enough that even when included in std dev calculation,
	// it still exceeds 3 * stdDev from centroid
	data := &models.SurveyData{
		ProjectID: "TEST",
		Points: []models.SurveyPoint{
			{PointID: "P1", Easting: 100, Northing: 100},
			{PointID: "P2", Easting: 100.1, Northing: 100},
			{PointID: "P3", Easting: 100, Northing: 100.1},
			{PointID: "P4", Easting: 100.1, Northing: 100.1},
			{PointID: "P5", Easting: 100.05, Northing: 100.05},
			{PointID: "P6", Easting: 100.02, Northing: 100.08},
			{PointID: "P7", Easting: 100.08, Northing: 100.02},
			{PointID: "P8", Easting: 100.03, Northing: 100.07},
			{PointID: "P9", Easting: 100.07, Northing: 100.03},
			{PointID: "P10", Easting: 100.04, Northing: 100.06},
			{PointID: "OUTLIER", Easting: 200, Northing: 200}, // Very far from tight cluster
		},
	}

	issues := DetectOutliers(data)

	// With many tightly clustered points and one far outlier, it should be detected
	found := false
	for _, issue := range issues {
		if len(issue.PointIDs) > 0 && issue.PointIDs[0] == "OUTLIER" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected OUTLIER point to be detected")
	}
}

func TestCalculateSummaryStatistics(t *testing.T) {
	height1 := 100.0
	data := &models.SurveyData{
		ProjectID: "TEST",
		Points: []models.SurveyPoint{
			{PointID: "T1", Easting: 100, Northing: 100, SurveyType: models.SurveyTypeTraverse, Height: &height1},
			{PointID: "T2", Easting: 200, Northing: 200, SurveyType: models.SurveyTypeTraverse},
			{PointID: "C1", Easting: 150, Northing: 150, SurveyType: models.SurveyTypeControl},
			{PointID: "D1", Easting: 175, Northing: 175, SurveyType: models.SurveyTypeDetail},
		},
	}

	stats := CalculateSummaryStatistics(data)

	if stats.TotalPoints != 4 {
		t.Errorf("TotalPoints = %d, expected 4", stats.TotalPoints)
	}
	if stats.TraversePoints != 2 {
		t.Errorf("TraversePoints = %d, expected 2", stats.TraversePoints)
	}
	if stats.ControlPoints != 1 {
		t.Errorf("ControlPoints = %d, expected 1", stats.ControlPoints)
	}
	if stats.DetailPoints != 1 {
		t.Errorf("DetailPoints = %d, expected 1", stats.DetailPoints)
	}
	if stats.PointsWithHeight != 1 {
		t.Errorf("PointsWithHeight = %d, expected 1", stats.PointsWithHeight)
	}
}
