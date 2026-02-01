package domain

import (
	"math"
	"testing"

	"github.com/survey-validator/models"
)

func TestDistance(t *testing.T) {
	tests := []struct {
		name     string
		p1       models.SurveyPoint
		p2       models.SurveyPoint
		expected float64
	}{
		{
			name:     "Same point",
			p1:       models.SurveyPoint{Easting: 100, Northing: 100},
			p2:       models.SurveyPoint{Easting: 100, Northing: 100},
			expected: 0,
		},
		{
			name:     "Horizontal distance",
			p1:       models.SurveyPoint{Easting: 100, Northing: 100},
			p2:       models.SurveyPoint{Easting: 200, Northing: 100},
			expected: 100,
		},
		{
			name:     "Vertical distance",
			p1:       models.SurveyPoint{Easting: 100, Northing: 100},
			p2:       models.SurveyPoint{Easting: 100, Northing: 200},
			expected: 100,
		},
		{
			name:     "Diagonal distance",
			p1:       models.SurveyPoint{Easting: 0, Northing: 0},
			p2:       models.SurveyPoint{Easting: 3, Northing: 4},
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Distance(&tt.p1, &tt.p2)
			if math.Abs(result-tt.expected) > 0.0001 {
				t.Errorf("Distance() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestBearing(t *testing.T) {
	tests := []struct {
		name     string
		p1       models.SurveyPoint
		p2       models.SurveyPoint
		expected float64
	}{
		{
			name:     "Due North",
			p1:       models.SurveyPoint{Easting: 100, Northing: 100},
			p2:       models.SurveyPoint{Easting: 100, Northing: 200},
			expected: 0,
		},
		{
			name:     "Due East",
			p1:       models.SurveyPoint{Easting: 100, Northing: 100},
			p2:       models.SurveyPoint{Easting: 200, Northing: 100},
			expected: 90,
		},
		{
			name:     "Due South",
			p1:       models.SurveyPoint{Easting: 100, Northing: 200},
			p2:       models.SurveyPoint{Easting: 100, Northing: 100},
			expected: 180,
		},
		{
			name:     "Due West",
			p1:       models.SurveyPoint{Easting: 200, Northing: 100},
			p2:       models.SurveyPoint{Easting: 100, Northing: 100},
			expected: 270,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Bearing(&tt.p1, &tt.p2)
			if math.Abs(result-tt.expected) > 0.0001 {
				t.Errorf("Bearing() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestCentroid(t *testing.T) {
	points := []models.SurveyPoint{
		{Easting: 0, Northing: 0},
		{Easting: 100, Northing: 0},
		{Easting: 100, Northing: 100},
		{Easting: 0, Northing: 100},
	}

	e, n := Centroid(points)

	if math.Abs(e-50) > 0.0001 {
		t.Errorf("Centroid easting = %v, expected 50", e)
	}
	if math.Abs(n-50) > 0.0001 {
		t.Errorf("Centroid northing = %v, expected 50", n)
	}
}

func TestBoundingBox(t *testing.T) {
	points := []models.SurveyPoint{
		{Easting: 10, Northing: 20},
		{Easting: 100, Northing: 50},
		{Easting: 50, Northing: 200},
		{Easting: 30, Northing: 80},
	}

	bbox := BoundingBox(points)

	if bbox.MinEasting != 10 {
		t.Errorf("MinEasting = %v, expected 10", bbox.MinEasting)
	}
	if bbox.MaxEasting != 100 {
		t.Errorf("MaxEasting = %v, expected 100", bbox.MaxEasting)
	}
	if bbox.MinNorthing != 20 {
		t.Errorf("MinNorthing = %v, expected 20", bbox.MinNorthing)
	}
	if bbox.MaxNorthing != 200 {
		t.Errorf("MaxNorthing = %v, expected 200", bbox.MaxNorthing)
	}
}
