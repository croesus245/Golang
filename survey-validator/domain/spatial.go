package domain

import (
	"math"

	"github.com/survey-validator/models"
)

func Distance(p1, p2 *models.SurveyPoint) float64 {
	dE := p2.Easting - p1.Easting
	dN := p2.Northing - p1.Northing
	return math.Sqrt(dE*dE + dN*dN)
}

func Distance3D(p1, p2 *models.SurveyPoint) float64 {
	if !p1.HasHeight() || !p2.HasHeight() {
		return Distance(p1, p2)
	}
	dE := p2.Easting - p1.Easting
	dN := p2.Northing - p1.Northing
	dH := *p2.Height - *p1.Height
	return math.Sqrt(dE*dE + dN*dN + dH*dH)
}

func Bearing(p1, p2 *models.SurveyPoint) float64 {
	dE := p2.Easting - p1.Easting
	dN := p2.Northing - p1.Northing
	bearing := math.Atan2(dE, dN) * 180 / math.Pi
	if bearing < 0 {
		bearing += 360
	}
	return bearing
}

func BearingDifference(b1, b2 float64) float64 {
	diff := math.Abs(b1 - b2)
	if diff > 180 {
		diff = 360 - diff
	}
	return diff
}

func Centroid(points []models.SurveyPoint) (float64, float64) {
	if len(points) == 0 {
		return 0, 0
	}
	var sumE, sumN float64
	for _, p := range points {
		sumE += p.Easting
		sumN += p.Northing
	}
	n := float64(len(points))
	return sumE / n, sumN / n
}

func StandardDeviation(points []models.SurveyPoint, cE, cN float64) float64 {
	if len(points) == 0 {
		return 0
	}
	var sum float64
	for _, p := range points {
		dE := p.Easting - cE
		dN := p.Northing - cN
		sum += dE*dE + dN*dN
	}
	return math.Sqrt(sum / float64(len(points)))
}

func BoundingBox(points []models.SurveyPoint) models.BBox {
	if len(points) == 0 {
		return models.BBox{}
	}
	bbox := models.BBox{
		MinEasting:  points[0].Easting,
		MaxEasting:  points[0].Easting,
		MinNorthing: points[0].Northing,
		MaxNorthing: points[0].Northing,
	}
	for _, p := range points[1:] {
		if p.Easting < bbox.MinEasting {
			bbox.MinEasting = p.Easting
		}
		if p.Easting > bbox.MaxEasting {
			bbox.MaxEasting = p.Easting
		}
		if p.Northing < bbox.MinNorthing {
			bbox.MinNorthing = p.Northing
		}
		if p.Northing > bbox.MaxNorthing {
			bbox.MaxNorthing = p.Northing
		}
	}
	return bbox
}
