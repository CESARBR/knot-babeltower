package entities

// Config represents the thing's config
type Config struct {
	SensorID int    `json:"sensorId"`
	Schema   Schema `json:"schema,omitempty"`
	Event    Event  `json:"event,omitempty"`
}
