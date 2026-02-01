package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/survey-validator/models"
)

// ValidationRequest represents the request body for validation
type ValidationRequest struct {
	ProjectID        string               `json:"project_id"`
	CoordinateSystem string               `json:"coordinate_system,omitempty"`
	Points           []models.SurveyPoint `json:"points"`
}

// ValidationResponse represents the response from validation
type ValidationResponse struct {
	*models.ValidationReport
}

// ValidateRequest validates the incoming request
func ValidateRequest(r *http.Request) (*models.SurveyData, error) {
	if r.Method != http.MethodPost {
		return nil, fmt.Errorf("method not allowed: %s", r.Method)
	}

	if r.Body == nil {
		return nil, fmt.Errorf("request body is required")
	}

	var req ValidationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	if len(req.Points) == 0 {
		return nil, fmt.Errorf("at least one point is required")
	}

	return &models.SurveyData{
		ProjectID:        req.ProjectID,
		CoordinateSystem: req.CoordinateSystem,
		Points:           req.Points,
	}, nil
}
