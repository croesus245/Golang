package domain

import (
	"testing"

	"github.com/survey-validator/models"
)

func TestComputeTraverseAdjustment_ClosedLoop(t *testing.T) {
	// simple closed traverse - square loop with small misclosure
	data := &models.SurveyData{
		ProjectID: "TEST-001",
		Points: []models.SurveyPoint{
			{PointID: "A", Easting: 1000.000, Northing: 1000.000, SurveyType: models.SurveyTypeTraverse},
			{PointID: "B", Easting: 1100.005, Northing: 1000.002, SurveyType: models.SurveyTypeTraverse},
			{PointID: "C", Easting: 1100.008, Northing: 1100.004, SurveyType: models.SurveyTypeTraverse},
			{PointID: "D", Easting: 1000.003, Northing: 1100.006, SurveyType: models.SurveyTypeTraverse},
			{PointID: "A", Easting: 1000.012, Northing: 1000.008, SurveyType: models.SurveyTypeTraverse}, // back to A with misclosure
		},
	}

	result := ComputeTraverseAdjustment(data, nil)

	if result.Status == "ERROR" {
		t.Fatalf("unexpected error: %s", result.Message)
	}

	// should have 4 legs
	if len(result.Legs) != 4 {
		t.Errorf("expected 4 legs, got %d", len(result.Legs))
	}

	// total distance should be about 400m
	if result.TotalDistance < 395 || result.TotalDistance > 405 {
		t.Errorf("expected total distance ~400m, got %.3f", result.TotalDistance)
	}

	// linear misclosure should be small
	if result.LinearMisclosure > 0.02 {
		t.Errorf("expected small misclosure, got %.4f", result.LinearMisclosure)
	}

	// precision should be good (>10000)
	if result.Precision < 10000 {
		t.Errorf("expected precision >1:10000, got 1:%.0f", result.Precision)
	}

	// should pass
	if result.Status != "PASS" {
		t.Errorf("expected PASS, got %s: %s", result.Status, result.Message)
	}

	t.Logf("Misclosure: %.4fm, Precision: %s", result.LinearMisclosure, result.ClosureRatio)
}

func TestComputeTraverseAdjustment_PoorClosure(t *testing.T) {
	// traverse with poor closure
	data := &models.SurveyData{
		ProjectID: "TEST-002",
		Points: []models.SurveyPoint{
			{PointID: "A", Easting: 1000.000, Northing: 1000.000, SurveyType: models.SurveyTypeTraverse},
			{PointID: "B", Easting: 1100.000, Northing: 1000.000, SurveyType: models.SurveyTypeTraverse},
			{PointID: "C", Easting: 1100.000, Northing: 1100.000, SurveyType: models.SurveyTypeTraverse},
			{PointID: "D", Easting: 1000.000, Northing: 1100.000, SurveyType: models.SurveyTypeTraverse},
			{PointID: "A", Easting: 1000.200, Northing: 1000.150, SurveyType: models.SurveyTypeTraverse}, // 25cm misclosure
		},
	}

	input := &models.TraverseInput{
		RequiredPrecision: 10000, // require 1:10000
	}

	result := ComputeTraverseAdjustment(data, input)

	// precision should be poor (~1:1600)
	if result.Precision > 2000 {
		t.Errorf("expected poor precision, got 1:%.0f", result.Precision)
	}

	// should fail
	if result.Status != "FAIL" {
		t.Errorf("expected FAIL, got %s", result.Status)
	}

	t.Logf("Misclosure: %.4fm, Precision: %s, Status: %s",
		result.LinearMisclosure, result.ClosureRatio, result.Status)
}

func TestComputeTraverseAdjustment_TooFewPoints(t *testing.T) {
	data := &models.SurveyData{
		ProjectID: "TEST-003",
		Points: []models.SurveyPoint{
			{PointID: "A", Easting: 1000.0, Northing: 1000.0, SurveyType: models.SurveyTypeTraverse},
			{PointID: "B", Easting: 1100.0, Northing: 1000.0, SurveyType: models.SurveyTypeTraverse},
		},
	}

	result := ComputeTraverseAdjustment(data, nil)

	if result.Status != "ERROR" {
		t.Errorf("expected ERROR for too few points, got %s", result.Status)
	}
}

func TestComputeTraverseAdjustment_AdjustedCoords(t *testing.T) {
	// check that adjusted coordinates are computed correctly
	data := &models.SurveyData{
		ProjectID: "TEST-004",
		Points: []models.SurveyPoint{
			{PointID: "A", Easting: 1000.000, Northing: 1000.000, SurveyType: models.SurveyTypeTraverse},
			{PointID: "B", Easting: 1100.010, Northing: 1000.000, SurveyType: models.SurveyTypeTraverse},
			{PointID: "C", Easting: 1100.010, Northing: 1100.010, SurveyType: models.SurveyTypeTraverse},
			{PointID: "A", Easting: 1000.000, Northing: 1100.010, SurveyType: models.SurveyTypeTraverse},
		},
	}

	result := ComputeTraverseAdjustment(data, nil)

	// should have adjusted points (A is held, B and C adjusted)
	if len(result.AdjustedPoints) < 3 {
		t.Errorf("expected at least 3 adjusted points, got %d", len(result.AdjustedPoints))
	}

	// first point should be unchanged
	if result.AdjustedPoints[0].AdjEasting != 1000.000 {
		t.Errorf("first point should be held fixed")
	}

	t.Logf("Adjusted points: %+v", result.AdjustedPoints)
}
