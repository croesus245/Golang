package engine

import (
	"testing"

	"github.com/survey-validator/models"
)

func TestEngine_Validate(t *testing.T) {
	engine := NewEngine()

	data := &models.SurveyData{
		ProjectID: "TEST-001",
		Points: []models.SurveyPoint{
			{PointID: "T1", Easting: 500000, Northing: 6000000, SurveyType: models.SurveyTypeTraverse},
			{PointID: "T2", Easting: 500050, Northing: 6000025, SurveyType: models.SurveyTypeTraverse},
			{PointID: "T3", Easting: 500100, Northing: 6000050, SurveyType: models.SurveyTypeTraverse},
		},
	}

	report := engine.Validate(data)

	if report.ProjectID != "TEST-001" {
		t.Errorf("ProjectID = %s, expected TEST-001", report.ProjectID)
	}

	if report.Summary.TotalPoints != 3 {
		t.Errorf("TotalPoints = %d, expected 3", report.Summary.TotalPoints)
	}

	if len(report.ChecksPerformed) == 0 {
		t.Error("Expected checks to be performed")
	}

	if report.ProcessingTime == "" {
		t.Error("Expected processing time to be recorded")
	}
}

func TestEngine_ValidateWithErrors(t *testing.T) {
	engine := NewEngine()

	data := &models.SurveyData{
		ProjectID: "TEST-002",
		Points: []models.SurveyPoint{
			{PointID: "P1", Easting: 100, Northing: 100},
			{PointID: "P2", Easting: 100.0001, Northing: 100.0001}, // Duplicate
		},
	}

	report := engine.Validate(data)

	if report.Status == models.StatusPass {
		t.Error("Expected status to not be PASS due to duplicates")
	}

	if len(report.Issues) == 0 {
		t.Error("Expected issues to be found")
	}
}

func TestEngine_EmptyData(t *testing.T) {
	engine := NewEngine()

	data := &models.SurveyData{
		ProjectID: "TEST-003",
		Points:    []models.SurveyPoint{},
	}

	report := engine.Validate(data)

	if report.Status != models.StatusFail {
		t.Errorf("Status = %s, expected FAIL for empty data", report.Status)
	}
}
