package models

// control.go - control extension mode structures

// NetworkType - type of control network being validated
type NetworkType string

const (
	NetworkTraverse NetworkType = "traverse"
	NetworkLeveling NetworkType = "leveling"
	NetworkGNSS     NetworkType = "gnss"
)

// TraverseType - classification of traverse
type TraverseType string

const (
	TraverseClosed TraverseType = "closed" // returns to start point
	TraverseLink   TraverseType = "link"   // known start AND end
	TraverseOpen   TraverseType = "open"   // only start known (weak)
)

// ToleranceClass - survey accuracy requirements
type ToleranceClass string

const (
	ClassFirstOrder   ToleranceClass = "first_order"  // 1:25000+
	ClassSecondOrder  ToleranceClass = "second_order" // 1:10000
	ClassThirdOrder   ToleranceClass = "third_order"  // 1:5000
	ClassEngineering  ToleranceClass = "engineering"  // 1:3000
	ClassConstruction ToleranceClass = "construction" // 1:1000
)

// tolerance requirements for each class
var TolerancePrecision = map[ToleranceClass]float64{
	ClassFirstOrder:   25000,
	ClassSecondOrder:  10000,
	ClassThirdOrder:   5000,
	ClassEngineering:  3000,
	ClassConstruction: 1000,
}

// KnownControl - a known control point with fixed coords
type KnownControl struct {
	PointID  string  `json:"point_id"`
	Easting  float64 `json:"easting"`
	Northing float64 `json:"northing"`
	Height   float64 `json:"height,omitempty"`
	IsStart  bool    `json:"is_start"`
	IsEnd    bool    `json:"is_end"`
}

// TraverseObservation - distance + angle observation (not coords)
type TraverseObservation struct {
	StationID string  `json:"station_id"`
	TargetID  string  `json:"target_id"`
	Distance  float64 `json:"distance"`             // horizontal distance in meters
	Bearing   float64 `json:"bearing,omitempty"`    // azimuth/bearing in degrees
	Angle     float64 `json:"angle,omitempty"`      // horizontal angle in degrees
	AngleType string  `json:"angle_type,omitempty"` // "left" or "right"
}

// LevelingObservation - single leveling reading
type LevelingObservation struct {
	PointID  string  `json:"point_id"`
	BS       float64 `json:"backsight,omitempty"`    // backsight reading
	IS       float64 `json:"intermediate,omitempty"` // intermediate sight
	FS       float64 `json:"foresight,omitempty"`    // foresight reading
	Distance float64 `json:"distance,omitempty"`     // distance for this setup (for allowable calc)
}

// ControlExtensionInput - full input for control extension mode
type ControlExtensionInput struct {
	Mode           string         `json:"mode"` // "topo" or "control_extension"
	NetworkType    NetworkType    `json:"network_type"`
	ToleranceClass ToleranceClass `json:"tolerance_class"`

	// known controls
	StartControl *KnownControl `json:"start_control,omitempty"`
	EndControl   *KnownControl `json:"end_control,omitempty"`

	// for observation-based input
	StartBearing float64               `json:"start_bearing,omitempty"`
	Observations []TraverseObservation `json:"observations,omitempty"`

	// for leveling
	StartBMHeight float64               `json:"start_bm_height,omitempty"`
	LevelingObs   []LevelingObservation `json:"leveling_obs,omitempty"`
}

// LevelingResult - computed leveling results
type LevelingResult struct {
	StartBM          string          `json:"start_bm"`
	StartHeight      float64         `json:"start_height"`
	EndBM            string          `json:"end_bm,omitempty"`
	EndHeight        float64         `json:"end_height,omitempty"`
	TotalDistance    float64         `json:"total_distance_km"`
	HeightMisclosure float64         `json:"height_misclosure"`
	AllowableMisc    float64         `json:"allowable_misclosure"`
	Status           string          `json:"status"`
	Points           []LevelingPoint `json:"points"`
	Message          string          `json:"message"`
}

// LevelingPoint - computed RL for each point
type LevelingPoint struct {
	PointID    string  `json:"point_id"`
	Rise       float64 `json:"rise,omitempty"`
	Fall       float64 `json:"fall,omitempty"`
	RawRL      float64 `json:"raw_rl"`
	AdjustedRL float64 `json:"adjusted_rl"`
	Correction float64 `json:"correction"`
}
