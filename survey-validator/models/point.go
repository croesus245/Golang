package models

type SurveyType string

const (
	SurveyTypeTraverse SurveyType = "traverse"
	SurveyTypeControl  SurveyType = "control"
	SurveyTypeDetail   SurveyType = "detail"
)

type SurveyPoint struct {
	PointID          string     `json:"point_id"`
	Easting          float64    `json:"easting"`
	Northing         float64    `json:"northing"`
	Height           *float64   `json:"height,omitempty"`
	SurveyType       SurveyType `json:"survey_type"`
	CoordinateSystem string     `json:"coordinate_system,omitempty"`
}

type SurveyData struct {
	ProjectID        string        `json:"project_id"`
	CoordinateSystem string        `json:"coordinate_system,omitempty"`
	Points           []SurveyPoint `json:"points"`
}

func (p *SurveyPoint) IsValid() bool {
	return p.PointID != "" && p.Easting != 0 && p.Northing != 0
}

func (p *SurveyPoint) HasHeight() bool {
	return p.Height != nil
}
