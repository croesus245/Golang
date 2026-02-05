package models

// report.go - validation results and stats

import "time"

type ValidationStatus string

const (
	StatusPass    ValidationStatus = "PASS"
	StatusWarning ValidationStatus = "WARNING"
	StatusFail    ValidationStatus = "FAIL"
)

type IssueSeverity string

const (
	SeverityError   IssueSeverity = "error"
	SeverityWarning IssueSeverity = "warning"
	SeverityInfo    IssueSeverity = "info"
)

type ValidationIssue struct {
	CheckName   string        `json:"check_name"`
	Severity    IssueSeverity `json:"severity"`
	PointIDs    []string      `json:"point_ids,omitempty"`
	Description string        `json:"description"`
	Details     interface{}   `json:"details,omitempty"`
}

type SummaryStatistics struct {
	TotalPoints      int     `json:"total_points"`
	TraversePoints   int     `json:"traverse_points"`
	ControlPoints    int     `json:"control_points"`
	DetailPoints     int     `json:"detail_points"`
	PointsWithHeight int     `json:"points_with_height"`
	BoundingBox      BBox    `json:"bounding_box"`
	CentroidEasting  float64 `json:"centroid_easting"`
	CentroidNorthing float64 `json:"centroid_northing"`
}

type BBox struct {
	MinEasting  float64 `json:"min_easting"`
	MaxEasting  float64 `json:"max_easting"`
	MinNorthing float64 `json:"min_northing"`
	MaxNorthing float64 `json:"max_northing"`
}

type TraverseClosureDetails struct {
	MisclosureEasting  float64 `json:"misclosure_easting"`
	MisclosureNorthing float64 `json:"misclosure_northing"`
	LinearMisclosure   float64 `json:"linear_misclosure"`
	TraverseLength     float64 `json:"traverse_length"`
	RelativePrecision  string  `json:"relative_precision"`
	Quality            string  `json:"quality"`
}

type ValidationReport struct {
	ProjectID       string            `json:"project_id"`
	Timestamp       time.Time         `json:"timestamp"`
	Status          ValidationStatus  `json:"status"`
	ConfidenceScore float64           `json:"confidence_score"`
	Summary         SummaryStatistics `json:"summary"`
	Issues          []ValidationIssue `json:"issues"`
	ChecksPerformed []string          `json:"checks_performed"`
	ProcessingTime  string            `json:"processing_time"`
	TraverseResult  *TraverseResult   `json:"traverse_adjustment,omitempty"`
}

// NewValidationReport - starts with PASS, we'll downgrade if issues found
func NewValidationReport(projectID string) *ValidationReport {
	return &ValidationReport{
		ProjectID: projectID,
		Timestamp: time.Now(),
		Status:    StatusPass,
		Issues:    make([]ValidationIssue, 0),
	}
}

// AddIssue and update status accordingly
func (r *ValidationReport) AddIssue(issue ValidationIssue) {
	r.Issues = append(r.Issues, issue)

	switch issue.Severity {
	case SeverityError:
		r.Status = StatusFail
	case SeverityWarning:
		if r.Status != StatusFail {
			r.Status = StatusWarning
		}
	}
}

// CalculateConfidenceScore - rough quality indicator
// starts at 100, deduct points for each issue
func (r *ValidationReport) CalculateConfidenceScore() {
	if len(r.Issues) == 0 {
		r.ConfidenceScore = 100.0
		return
	}

	score := 100.0
	for _, issue := range r.Issues {
		switch issue.Severity {
		case SeverityError:
			score -= 15.0
		case SeverityWarning:
			score -= 5.0
		case SeverityInfo:
			score -= 1.0
		}
	}

	if score < 0 {
		score = 0
	}
	r.ConfidenceScore = score
}
