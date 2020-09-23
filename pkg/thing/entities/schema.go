package entities

// Schema represents the thing's schema
type Schema struct {
	ValueType int    `json:"valueType" validate:"required"`
	Unit      int    `json:"unit"`
	TypeID    int    `json:"typeId" validate:"required"`
	Name      string `json:"name" validate:"required,max=30"`
}
