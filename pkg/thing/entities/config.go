package entities

// Config represents the thing's config
type Config struct {
	SensorID       int         `json:"sensorId"`
	Change         bool        `json:"change"`
	TimeSec        int         `json:"timeSec,omitempty"`
	LowerThreshold interface{} `json:"lowerThreshold,omitempty"`
	UpperThreshold interface{} `json:"upperThreshold,omitempty"`
}
