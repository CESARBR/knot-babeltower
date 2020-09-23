package entities

// Event represents the thing's event
type Event struct {
	Change         bool        `json:"change"`
	TimeSec        int         `json:"timeSec,omitempty"`
	LowerThreshold interface{} `json:"lowerThreshold,omitempty"`
	UpperThreshold interface{} `json:"upperThreshold,omitempty"`
}
