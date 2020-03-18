package entities

// Data represents the thing's data
type Data struct {
	SensorID int         `json:"sensorId"`
	Value    interface{} `json:"value"`
}
