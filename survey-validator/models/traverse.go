package models

// traverse.go - structures for traverse computations and adjustments

// TraverseLeg - one leg of the traverse with computed values
type TraverseLeg struct {
	FromPoint   string  `json:"from_point"`
	ToPoint     string  `json:"to_point"`
	Distance    float64 `json:"distance"`
	Bearing     float64 `json:"bearing"`
	DeltaE      float64 `json:"delta_e"`
	DeltaN      float64 `json:"delta_n"`
	CorrectionE float64 `json:"correction_e"`
	CorrectionN float64 `json:"correction_n"`
	AdjustedDE  float64 `json:"adjusted_delta_e"`
	AdjustedDN  float64 `json:"adjusted_delta_n"`
}

// AdjustedPoint - final adjusted coordinates for a traverse station
type AdjustedPoint struct {
	PointID      string  `json:"point_id"`
	RawEasting   float64 `json:"raw_easting"`
	RawNorthing  float64 `json:"raw_northing"`
	AdjEasting   float64 `json:"adjusted_easting"`
	AdjNorthing  float64 `json:"adjusted_northing"`
	ResidualE    float64 `json:"residual_e"`
	ResidualN    float64 `json:"residual_n"`
	ResidualDist float64 `json:"residual_distance"`
}

// TraverseResult - complete adjustment output
type TraverseResult struct {
	// Classification
	TraverseType     string `json:"traverse_type"` // closed, link, open
	TraverseTypeDesc string `json:"traverse_type_desc"`

	// Misclosure info
	SumDeltaE        float64 `json:"sum_delta_e"`
	SumDeltaN        float64 `json:"sum_delta_n"`
	LinearMisclosure float64 `json:"linear_misclosure"`
	TotalDistance    float64 `json:"total_distance"`
	ClosureRatio     string  `json:"closure_ratio"`
	Precision        float64 `json:"precision"`

	// Angular (if provided)
	AngularMisclosure float64 `json:"angular_misclosure,omitempty"`
	AllowableAngular  float64 `json:"allowable_angular,omitempty"`
	AngularStatus     string  `json:"angular_status,omitempty"`

	// Results
	Legs           []TraverseLeg   `json:"legs"`
	AdjustedPoints []AdjustedPoint `json:"adjusted_points"`

	// Pass/Fail
	Status            string   `json:"status"`
	RequiredPrecision float64  `json:"required_precision"`
	Message           string   `json:"message"`
	SuggestedFixes    []string `json:"suggested_fixes,omitempty"`
}

// TraverseInput - optional extended input for angle-based traverses
type TraverseInput struct {
	Stations          []TraverseStation `json:"stations,omitempty"`
	StartBearing      float64           `json:"start_bearing,omitempty"`
	RequiredPrecision float64           `json:"required_precision,omitempty"` // e.g. 5000 for 1:5000
}

// TraverseStation - for angle/distance input method
type TraverseStation struct {
	PointID  string  `json:"point_id"`
	Angle    float64 `json:"angle,omitempty"`    // horizontal angle at this station
	Distance float64 `json:"distance,omitempty"` // distance to next station
}
