package domain

// leveling.go - level run computations for vertical control extension

import (
	"fmt"
	"math"

	"github.com/survey-validator/models"
)

// allowable misclosure constants (mm per sqrt(km))
const (
	FirstOrderLeveling  = 3.0  // 3mm√K
	SecondOrderLeveling = 6.0  // 6mm√K
	ThirdOrderLeveling  = 12.0 // 12mm√K
	EngineeringLeveling = 24.0 // 24mm√K
)

// ComputeLeveling - processes level run observations and returns adjusted RLs
func ComputeLeveling(obs []models.LevelingObservation, startHeight float64, endHeight float64, toleranceClass models.ToleranceClass) *models.LevelingResult {
	result := &models.LevelingResult{
		StartHeight: startHeight,
		Points:      make([]models.LevelingPoint, 0),
	}

	if len(obs) < 2 {
		result.Status = "ERROR"
		result.Message = "Need at least 2 observations for leveling"
		return result
	}

	result.StartBM = obs[0].PointID

	// compute rise/fall and raw RLs using rise and fall method
	currentRL := startHeight
	var totalDistance float64
	var sumRise, sumFall float64

	for i, o := range obs {
		pt := models.LevelingPoint{
			PointID: o.PointID,
			RawRL:   currentRL,
		}

		totalDistance += o.Distance

		if i > 0 {
			prevObs := obs[i-1]

			// if we have BS and FS, compute rise/fall
			if prevObs.FS > 0 && o.BS > 0 {
				// this is a change point or turning point
				// rise = BS - FS (previous setup)
				// but for simple case: rise = BS(prev) - FS(curr)
			}

			// simple rise/fall from consecutive readings
			// rise = BS - FS (positive = ground went up)
			if o.FS > 0 && prevObs.BS > 0 {
				diff := prevObs.BS - o.FS
				if diff > 0 {
					pt.Rise = diff
					sumRise += diff
				} else {
					pt.Fall = -diff
					sumFall += diff
				}
				currentRL = currentRL + diff
				pt.RawRL = currentRL
			} else if o.IS > 0 && prevObs.BS > 0 {
				// intermediate sight
				diff := prevObs.BS - o.IS
				currentRL = result.StartHeight + sumRise - sumFall + diff
				pt.RawRL = currentRL
			}
		}

		result.Points = append(result.Points, pt)
	}

	// last point
	lastPt := obs[len(obs)-1]
	result.EndBM = lastPt.PointID
	result.EndHeight = currentRL
	result.TotalDistance = totalDistance / 1000.0 // convert to km

	// compute misclosure if end height is known
	if endHeight > 0 {
		result.HeightMisclosure = currentRL - endHeight
	} else {
		// closed level loop - should return to start height
		result.HeightMisclosure = currentRL - startHeight
	}

	// allowable misclosure based on class
	c := getAllowableConstant(toleranceClass)
	if result.TotalDistance > 0 {
		result.AllowableMisc = c * math.Sqrt(result.TotalDistance) / 1000.0 // convert to meters
	} else {
		result.AllowableMisc = c / 1000.0 // assume 1km if no distance
	}

	// pass/fail
	if math.Abs(result.HeightMisclosure) <= result.AllowableMisc {
		result.Status = "PASS"
		result.Message = fmt.Sprintf("Level run acceptable: %.4fm misclosure within %.4fm allowable",
			result.HeightMisclosure, result.AllowableMisc)
	} else {
		result.Status = "FAIL"
		result.Message = fmt.Sprintf("Level run FAILED: %.4fm misclosure exceeds %.4fm allowable",
			result.HeightMisclosure, result.AllowableMisc)
	}

	// apply corrections (proportional to distance or equal per setup)
	applyLevelingCorrections(result)

	return result
}

// getAllowableConstant returns mm per sqrt(km) for each class
func getAllowableConstant(class models.ToleranceClass) float64 {
	switch class {
	case models.ClassFirstOrder:
		return FirstOrderLeveling
	case models.ClassSecondOrder:
		return SecondOrderLeveling
	case models.ClassThirdOrder:
		return ThirdOrderLeveling
	case models.ClassEngineering:
		return EngineeringLeveling
	default:
		return ThirdOrderLeveling // default
	}
}

// applyLevelingCorrections distributes misclosure across all points
func applyLevelingCorrections(result *models.LevelingResult) {
	if len(result.Points) < 2 {
		return
	}

	// distribute correction equally across points (simple method)
	// more sophisticated: proportional to distance
	corrPerPoint := -result.HeightMisclosure / float64(len(result.Points)-1)

	for i := range result.Points {
		if i == 0 {
			result.Points[i].AdjustedRL = result.Points[i].RawRL
			result.Points[i].Correction = 0
		} else {
			result.Points[i].Correction = corrPerPoint * float64(i)
			result.Points[i].AdjustedRL = math.Round((result.Points[i].RawRL+result.Points[i].Correction)*10000) / 10000
		}
	}
}

// ComputeLevelingFromRiseFall - alternative input: direct rise/fall values
func ComputeLevelingFromRiseFall(points []models.LevelingPoint, startHeight float64, endHeight float64) *models.LevelingResult {
	result := &models.LevelingResult{
		StartHeight: startHeight,
		Points:      make([]models.LevelingPoint, 0),
	}

	currentRL := startHeight

	for i, p := range points {
		pt := p
		if i == 0 {
			pt.RawRL = startHeight
		} else {
			currentRL = currentRL + p.Rise - p.Fall
			pt.RawRL = currentRL
		}
		result.Points = append(result.Points, pt)
	}

	result.EndHeight = currentRL

	if endHeight > 0 {
		result.HeightMisclosure = currentRL - endHeight
	}

	return result
}
