package domain

// adjustment.go - Bowditch traverse adjustment (compass rule)

import (
	"fmt"
	"math"

	"github.com/survey-validator/models"
)

// default tolerance if not specified
const DefaultRequiredPrecision = 5000.0 // 1:5000

// ComputeTraverseAdjustment - main function for closed traverse adjustment
// Takes traverse points in order, computes misclosure, applies Bowditch
func ComputeTraverseAdjustment(data *models.SurveyData, input *models.TraverseInput) *models.TraverseResult {
	result := &models.TraverseResult{
		Legs:           make([]models.TraverseLeg, 0),
		AdjustedPoints: make([]models.AdjustedPoint, 0),
		SuggestedFixes: make([]string, 0),
	}

	// filter to traverse points only, keep order
	var pts []models.SurveyPoint
	for _, p := range data.Points {
		if p.SurveyType == models.SurveyTypeTraverse {
			pts = append(pts, p)
		}
	}

	if len(pts) < 3 {
		result.Status = "ERROR"
		result.Message = "Need at least 3 traverse points for adjustment"
		return result
	}

	// set required precision
	result.RequiredPrecision = DefaultRequiredPrecision
	if input != nil && input.RequiredPrecision > 0 {
		result.RequiredPrecision = input.RequiredPrecision
	}

	// Step 1: compute legs (deltas) from coordinates
	var totalDist, sumDE, sumDN float64
	var longestLegIdx int
	var longestLegDist float64

	for i := 0; i < len(pts)-1; i++ {
		p1 := pts[i]
		p2 := pts[i+1]

		dE := p2.Easting - p1.Easting
		dN := p2.Northing - p1.Northing
		dist := math.Sqrt(dE*dE + dN*dN)
		bearing := calcBearing(dE, dN)

		leg := models.TraverseLeg{
			FromPoint: p1.PointID,
			ToPoint:   p2.PointID,
			Distance:  dist,
			Bearing:   bearing,
			DeltaE:    dE,
			DeltaN:    dN,
		}

		result.Legs = append(result.Legs, leg)
		totalDist += dist
		sumDE += dE
		sumDN += dN

		if dist > longestLegDist {
			longestLegDist = dist
			longestLegIdx = i
		}
	}

	// Step 2: classify traverse type
	first := pts[0]
	last := pts[len(pts)-1]
	closureDist := math.Sqrt(math.Pow(last.Easting-first.Easting, 2) + math.Pow(last.Northing-first.Northing, 2))

	// classify: closed, link, or open
	if first.PointID == last.PointID || closureDist < 0.001 {
		result.TraverseType = "closed"
		result.TraverseTypeDesc = "Closed traverse - returns to start point"
	} else if closureDist < totalDist*0.1 {
		// close enough to be considered a loop that should close
		result.TraverseType = "closed"
		result.TraverseTypeDesc = "Closed traverse - loop with misclosure"
	} else {
		// check if we have known end control (would make it a link traverse)
		// for now, treat as open unless very close to start
		result.TraverseType = "open"
		result.TraverseTypeDesc = "Open traverse - end point not at start (weak geometry)"
		result.SuggestedFixes = append(result.SuggestedFixes,
			"Consider closing the traverse back to start point for stronger geometry")
	}

	isClosed := result.TraverseType == "closed"

	if !isClosed {
		sumDE = last.Easting - first.Easting
		sumDN = last.Northing - first.Northing
	}

	// Step 3: linear misclosure
	result.SumDeltaE = sumDE
	result.SumDeltaN = sumDN
	result.LinearMisclosure = math.Sqrt(sumDE*sumDE + sumDN*sumDN)
	result.TotalDistance = totalDist

	// precision ratio
	if result.LinearMisclosure > 0 {
		result.Precision = totalDist / result.LinearMisclosure
		result.ClosureRatio = fmt.Sprintf("1:%.0f", result.Precision)
	} else {
		result.Precision = math.MaxFloat64
		result.ClosureRatio = "1:∞ (perfect)"
	}

	// Step 4: Bowditch adjustment (corrections proportional to distance)
	// correction per leg = (leg distance / total distance) * misclosure

	for i := range result.Legs {
		leg := &result.Legs[i]
		proportion := leg.Distance / totalDist

		// corrections are negative of misclosure, distributed by proportion
		leg.CorrectionE = -sumDE * proportion
		leg.CorrectionN = -sumDN * proportion
		leg.AdjustedDE = leg.DeltaE + leg.CorrectionE
		leg.AdjustedDN = leg.DeltaN + leg.CorrectionN
	}

	// Step 5: compute adjusted coordinates
	// start from first point (held fixed)
	cumE := first.Easting
	cumN := first.Northing

	// first point is fixed
	result.AdjustedPoints = append(result.AdjustedPoints, models.AdjustedPoint{
		PointID:     first.PointID,
		RawEasting:  first.Easting,
		RawNorthing: first.Northing,
		AdjEasting:  first.Easting,
		AdjNorthing: first.Northing,
	})

	// compute running adjusted coords
	for i, leg := range result.Legs {
		cumE += leg.AdjustedDE
		cumN += leg.AdjustedDN

		rawPt := pts[i+1]
		residE := cumE - rawPt.Easting
		residN := cumN - rawPt.Northing
		residDist := math.Sqrt(residE*residE + residN*residN)

		// don't add duplicate for closed traverse where last=first
		if rawPt.PointID != first.PointID {
			result.AdjustedPoints = append(result.AdjustedPoints, models.AdjustedPoint{
				PointID:      rawPt.PointID,
				RawEasting:   rawPt.Easting,
				RawNorthing:  rawPt.Northing,
				AdjEasting:   round3(cumE),
				AdjNorthing:  round3(cumN),
				ResidualE:    round4(residE),
				ResidualN:    round4(residN),
				ResidualDist: round4(residDist),
			})
		}
	}

	// Step 6: pass/fail check and suggested fixes
	if result.Precision >= result.RequiredPrecision {
		result.Status = "PASS"
		result.Message = fmt.Sprintf("Traverse meets 1:%.0f requirement (achieved 1:%.0f)",
			result.RequiredPrecision, result.Precision)
	} else {
		result.Status = "FAIL"
		result.Message = fmt.Sprintf("Traverse does NOT meet 1:%.0f requirement (achieved 1:%.0f)",
			result.RequiredPrecision, result.Precision)

		// add suggested fixes based on analysis
		result.SuggestedFixes = append(result.SuggestedFixes,
			generateSuggestedFixes(result, pts, longestLegIdx)...)
	}

	return result
}

// generateSuggestedFixes - analyzes traverse and suggests what to re-check
func generateSuggestedFixes(result *models.TraverseResult, pts []models.SurveyPoint, longestLegIdx int) []string {
	fixes := []string{}

	// suggest checking the leg that contributes most to error
	if len(result.Legs) > 0 {
		// find leg with largest correction (likely source of error)
		var maxCorrIdx int
		var maxCorr float64
		for i, leg := range result.Legs {
			corr := math.Sqrt(leg.CorrectionE*leg.CorrectionE + leg.CorrectionN*leg.CorrectionN)
			if corr > maxCorr {
				maxCorr = corr
				maxCorrIdx = i
			}
		}

		if maxCorr > 0.01 { // only suggest if correction is significant
			leg := result.Legs[maxCorrIdx]
			fixes = append(fixes,
				fmt.Sprintf("Re-check distance %s to %s (largest correction: %.4fm)",
					leg.FromPoint, leg.ToPoint, maxCorr))
		}

		// if easting error is much larger than northing, suggest bearing check
		if math.Abs(result.SumDeltaE) > math.Abs(result.SumDeltaN)*2 {
			fixes = append(fixes, "Easting error dominant - check angles/bearings for E-W pointing legs")
		} else if math.Abs(result.SumDeltaN) > math.Abs(result.SumDeltaE)*2 {
			fixes = append(fixes, "Northing error dominant - check angles/bearings for N-S pointing legs")
		}

		// check for short legs (more prone to angle errors)
		avgDist := result.TotalDistance / float64(len(result.Legs))
		for _, leg := range result.Legs {
			if leg.Distance < avgDist*0.3 {
				fixes = append(fixes,
					fmt.Sprintf("Short leg %s-%s (%.2fm) - angle errors have larger effect on short legs",
						leg.FromPoint, leg.ToPoint, leg.Distance))
				break // only suggest once
			}
		}
	}

	return fixes
}

// calcBearing - bearing from delta E/N
func calcBearing(dE, dN float64) float64 {
	bearing := math.Atan2(dE, dN) * 180 / math.Pi
	if bearing < 0 {
		bearing += 360
	}
	return bearing
}

// round helpers for cleaner output
func round3(v float64) float64 {
	return math.Round(v*1000) / 1000
}

func round4(v float64) float64 {
	return math.Round(v*10000) / 10000
}
