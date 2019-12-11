package entities

// Schema represents the thing's schema
type Schema struct {
	SensorID  int    `json:"sensorId"`
	ValueType int    `json:"valueType"`
	Unit      int    `json:"unit"`
	TypeID    int    `json:"typeId"`
	Name      string `json:"name"`
}
