package models

// point.go - core survey data structures

type SurveyType string

// survey types - what kind of point is this
const (
	SurveyTypeTraverse SurveyType = "traverse"
	SurveyTypeControl  SurveyType = "control"
	SurveyTypeDetail   SurveyType = "detail"
)

// SurveyPoint - a single measured point
type SurveyPoint struct {
	PointID          string     `json:"point_id"`
	Easting          float64    `json:"easting"`
	Northing         float64    `json:"northing"`
	Height           *float64   `json:"height,omitempty"`
	SurveyType       SurveyType `json:"survey_type"`
	CoordinateSystem string     `json:"coordinate_system,omitempty"`
}

// SurveyData - what comes in from the API
type SurveyData struct {
	ProjectID        string        `json:"project_id"`
	CoordinateSystem string        `json:"coordinate_system,omitempty"`
	Points           []SurveyPoint `json:"points"`
}

// IsValid - basic check, has id and non-zero coords
func (p *SurveyPoint) IsValid() bool {
	return p.PointID != "" && p.Easting != 0 && p.Northing != 0
}

func (p *SurveyPoint) HasHeight() bool {
	return p.Height != nil
}
