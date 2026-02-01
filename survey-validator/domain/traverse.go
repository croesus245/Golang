package domain

import (
	"fmt"
	"math"

	"github.com/survey-validator/models"
)

func CheckDistanceAndBearing(data *models.SurveyData) []models.ValidationIssue {
	var issues []models.ValidationIssue

	var traversePoints []models.SurveyPoint
	for _, p := range data.Points {
		if p.SurveyType == models.SurveyTypeTraverse {
			traversePoints = append(traversePoints, p)
		}
	}

	if len(traversePoints) < 2 {
		return issues
	}

	var prevBearing, prevDist float64

	for i := 1; i < len(traversePoints); i++ {
		p1 := &traversePoints[i-1]
		p2 := &traversePoints[i]
		dist := Distance(p1, p2)
		bearing := Bearing(p1, p2)

		if dist < MinTraverseDistance {
			msg := fmt.Sprintf("Very short distance between %s and %s: %.4fm",
				p1.PointID, p2.PointID, dist)
			issues = append(issues, models.ValidationIssue{
				CheckName:   "distance_bearing_check",
				Severity:    models.SeverityWarning,
				PointIDs:    []string{p1.PointID, p2.PointID},
				Description: msg,
				Details:     map[string]interface{}{"distance": dist},
			})
		}

		if i > 1 {
			change := BearingDifference(bearing, prevBearing)
			if change > MaxBearingChange {
				msg := fmt.Sprintf("Large bearing change at %s: %.1f°", p1.PointID, change)
				issues = append(issues, models.ValidationIssue{
					CheckName:   "distance_bearing_check",
					Severity:    models.SeverityWarning,
					PointIDs:    []string{p1.PointID, p2.PointID},
					Description: msg,
				})
			}

			if prevDist > 0 {
				ratio := dist / prevDist
				if ratio > 10 || ratio < 0.1 {
					msg := fmt.Sprintf("Unusual distance ratio at %s: %.1f", p1.PointID, ratio)
					issues = append(issues, models.ValidationIssue{
						CheckName:   "distance_bearing_check",
						Severity:    models.SeverityInfo,
						PointIDs:    []string{p1.PointID, p2.PointID},
						Description: msg,
					})
				}
			}
		}

		prevBearing = bearing
		prevDist = dist
	}
	return issues
}

func CheckTraverseClosure(data *models.SurveyData) []models.ValidationIssue {
	var issues []models.ValidationIssue

	var traversePoints []models.SurveyPoint
	for _, p := range data.Points {
		if p.SurveyType == models.SurveyTypeTraverse {
			traversePoints = append(traversePoints, p)
		}
	}

	if len(traversePoints) < 3 {
		return issues
	}

	first := &traversePoints[0]
	last := &traversePoints[len(traversePoints)-1]
	closureDist := Distance(first, last)

	var totalLen float64
	for i := 1; i < len(traversePoints); i++ {
		totalLen += Distance(&traversePoints[i-1], &traversePoints[i])
	}

	if closureDist >= totalLen*0.1 {
		return issues
	}

	miscE := last.Easting - first.Easting
	miscN := last.Northing - first.Northing
	linMisc := math.Sqrt(miscE*miscE + miscN*miscN)

	var precision float64
	if linMisc > 0 {
		precision = totalLen / linMisc
	} else {
		precision = math.MaxFloat64
	}

	var quality string
	var severity models.IssueSeverity

	switch {
	case precision >= GoodPrecision:
		quality = "Good (better than 1:10000)"
		severity = models.SeverityInfo
	case precision >= AcceptablePrecision:
		quality = "Acceptable (1:5000 to 1:10000)"
		severity = models.SeverityInfo
	case precision >= 1000:
		quality = "Poor (1:1000 to 1:5000)"
		severity = models.SeverityWarning
	default:
		quality = "Unacceptable (worse than 1:1000)"
		severity = models.SeverityError
	}

	msg := fmt.Sprintf("Traverse closure: %.4fm misclosure, 1:%.0f precision (%s)",
		linMisc, precision, quality)

	issues = append(issues, models.ValidationIssue{
		CheckName:   "traverse_closure",
		Severity:    severity,
		PointIDs:    []string{first.PointID, last.PointID},
		Description: msg,
		Details: models.TraverseClosureDetails{
			MisclosureEasting:  miscE,
			MisclosureNorthing: miscN,
			LinearMisclosure:   linMisc,
			TraverseLength:     totalLen,
			RelativePrecision:  fmt.Sprintf("1:%.0f", precision),
			Quality:            quality,
		},
	})
	return issues
}

func CalculateSummaryStatistics(data *models.SurveyData) models.SummaryStatistics {
	stats := models.SummaryStatistics{TotalPoints: len(data.Points)}

	if len(data.Points) == 0 {
		return stats
	}

	for _, p := range data.Points {
		switch p.SurveyType {
		case models.SurveyTypeTraverse:
			stats.TraversePoints++
		case models.SurveyTypeControl:
			stats.ControlPoints++
		case models.SurveyTypeDetail:
			stats.DetailPoints++
		}
		if p.HasHeight() {
			stats.PointsWithHeight++
		}
	}

	stats.BoundingBox = BoundingBox(data.Points)
	stats.CentroidEasting, stats.CentroidNorthing = Centroid(data.Points)
	return stats
}
