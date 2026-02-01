package domain

import (
	"fmt"

	"github.com/survey-validator/models"
)

const (
	DuplicateThreshold     = 0.001  // meters
	NearDuplicateThreshold = 0.01   // meters
	OutlierThreshold       = 3.0    // standard deviations
	MaxBearingChange       = 170.0  // degrees
	MinTraverseDistance    = 0.1    // meters
	GoodPrecision          = 10000  // 1:10000
	AcceptablePrecision    = 5000   // 1:5000
)

func ValidateInput(data *models.SurveyData) []models.ValidationIssue {
	var issues []models.ValidationIssue

	if len(data.Points) == 0 {
		issues = append(issues, models.ValidationIssue{
			CheckName:   "input_validation",
			Severity:    models.SeverityError,
			Description: "No survey points provided",
		})
		return issues
	}

	for _, p := range data.Points {
		if p.PointID == "" {
			issues = append(issues, models.ValidationIssue{
				CheckName:   "input_validation",
				Severity:    models.SeverityError,
				Description: "Point found with empty Point ID",
			})
		}

		if p.Easting == 0 && p.Northing == 0 {
			msg := fmt.Sprintf("Point %s has zero coordinates", p.PointID)
			issues = append(issues, models.ValidationIssue{
				CheckName:   "input_validation",
				Severity:    models.SeverityWarning,
				PointIDs:    []string{p.PointID},
				Description: msg,
			})
		}

		switch p.SurveyType {
		case models.SurveyTypeTraverse, models.SurveyTypeControl, models.SurveyTypeDetail:
			// valid type
		default:
			if p.SurveyType != "" {
				msg := fmt.Sprintf("Point %s has unknown type: %s", p.PointID, p.SurveyType)
				issues = append(issues, models.ValidationIssue{
					CheckName:   "input_validation",
					Severity:    models.SeverityWarning,
					PointIDs:    []string{p.PointID},
					Description: msg,
				})
			}
		}
	}
	return issues
}

func DetectDuplicates(data *models.SurveyData) []models.ValidationIssue {
	var issues []models.ValidationIssue
	points := data.Points

	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			dist := Distance(&points[i], &points[j])

			if dist < DuplicateThreshold {
				msg := fmt.Sprintf("Duplicate points: %s and %s (%.4fm apart)",
					points[i].PointID, points[j].PointID, dist)
				issues = append(issues, models.ValidationIssue{
					CheckName:   "duplicate_detection",
					Severity:    models.SeverityError,
					PointIDs:    []string{points[i].PointID, points[j].PointID},
					Description: msg,
					Details:     map[string]interface{}{"distance": dist},
				})
			} else if dist < NearDuplicateThreshold {
				msg := fmt.Sprintf("Near-duplicate points: %s and %s (%.4fm apart)",
					points[i].PointID, points[j].PointID, dist)
				issues = append(issues, models.ValidationIssue{
					CheckName:   "duplicate_detection",
					Severity:    models.SeverityWarning,
					PointIDs:    []string{points[i].PointID, points[j].PointID},
					Description: msg,
					Details:     map[string]interface{}{"distance": dist},
				})
			}
		}
	}
	return issues
}

func DetectOutliers(data *models.SurveyData) []models.ValidationIssue {
	var issues []models.ValidationIssue

	if len(data.Points) < 3 {
		return issues
	}

	cE, cN := Centroid(data.Points)
	stdDev := StandardDeviation(data.Points, cE, cN)
	if stdDev == 0 {
		return issues
	}

	threshold := OutlierThreshold * stdDev
	centroid := &models.SurveyPoint{Easting: cE, Northing: cN}

	for _, p := range data.Points {
		dist := Distance(centroid, &p)
		if dist > threshold {
			msg := fmt.Sprintf("Point %s may be an outlier (%.1fm from centroid)",
				p.PointID, dist)
			issues = append(issues, models.ValidationIssue{
				CheckName:   "outlier_detection",
				Severity:    models.SeverityWarning,
				PointIDs:    []string{p.PointID},
				Description: msg,
				Details: map[string]interface{}{
					"distance":  dist,
					"threshold": threshold,
				},
			})
		}
	}
	return issues
}
