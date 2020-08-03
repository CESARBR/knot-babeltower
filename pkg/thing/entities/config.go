package entities

// Config represents the thing's config
type Config struct {
	SensorID       int
	Change         bool
	TimeSec        int
	LowerThreshold int
	UpperThreshold int
}
