package interactors

import (
	"fmt"

	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
	"github.com/go-playground/validator"
)

// ErrInvalidSchema represents the error when the schema has a invalid format
type ErrInvalidSchema struct{}

func (eis *ErrInvalidSchema) Error() string {
	return "Thing's schema is invalid"
}

type interval struct {
	min int
	max int
}

type schemaType struct {
	valueType interface{}
	unit      interface{}
}

// UpdateSchema receive the new sensor schema and update it on the thing's service
func (i *ThingInteractor) UpdateSchema(authorization, thingID string, schemaList []entities.Schema) error {

	if !i.isValidSchema(schemaList) {
		return &ErrInvalidSchema{}
	}

	err := i.thingProxy.UpdateSchema(authorization, thingID, schemaList)
	if err != nil {
		return err
	}

	i.logger.Info("Schema updated")

	err = i.msgPublisher.SendUpdatedSchema(thingID)
	if err != nil {
		return err
	}

	return nil
}

func (i *ThingInteractor) isValidSchema(schemaList []entities.Schema) bool {
	validate := validator.New()
	validate.RegisterStructValidation(schemaValidation, entities.Schema{})
	for _, schema := range schemaList {
		err := validate.Struct(schema)
		if err != nil {
			fmt.Println(err)
			return false
		}
	}

	return true
}

var rules = map[int]schemaType{
	0:      schemaType{valueType: interval{1, 4}, unit: 0},
	0xfff0: schemaType{valueType: interval{1, 4}, unit: 0},
	0xfff1: schemaType{valueType: interval{1, 4}, unit: 0},
	0xfff2: schemaType{valueType: interval{1, 4}, unit: 0},
	0xff10: schemaType{valueType: interval{1, 4}, unit: 0},
	1:      schemaType{valueType: 1, unit: interval{1, 3}},
	2:      schemaType{valueType: 1, unit: interval{1, 2}},
	3:      schemaType{valueType: 1, unit: 1},
	4:      schemaType{valueType: 1, unit: interval{1, 3}},
	5:      schemaType{valueType: 1, unit: interval{1, 3}},
	6:      schemaType{valueType: 1, unit: 1},
	7:      schemaType{valueType: 1, unit: interval{1, 3}},
	8:      schemaType{valueType: 1, unit: interval{1, 3}},
	9:      schemaType{valueType: 1, unit: interval{1, 4}},
	0x0A:   schemaType{valueType: 1, unit: interval{1, 3}},
	0x0B:   schemaType{valueType: 1, unit: interval{1, 4}},
	0x0C:   schemaType{valueType: 2, unit: interval{1, 2}},
	0x0D:   schemaType{valueType: 2, unit: interval{1, 4}},
	0x0E:   schemaType{valueType: 2, unit: interval{1, 3}},
	0x0F:   schemaType{valueType: 2, unit: 1},
	0x10:   schemaType{valueType: 2, unit: 1},
	0x11:   schemaType{valueType: 2, unit: 1},
	0x12:   schemaType{valueType: 2, unit: 1},
	0x13:   schemaType{valueType: 1, unit: interval{1, 4}},
	0x14:   schemaType{valueType: 2, unit: interval{1, 6}},
	0x15:   schemaType{valueType: 1, unit: interval{1, 6}},
}

func isValidValueType(typeID, valueType int) bool {
	t := rules[typeID].valueType
	if t == nil {
		return false
	}

	switch v := t.(type) {
	case int:
		value := v
		if valueType != value {
			return false
		}
	case interval:
		interval := t.(interval)
		if valueType < interval.min || interval.max < valueType {
			return false
		}
	}

	return true
}

func isValidUnit(typeID, unit int) bool {
	u := rules[typeID].unit
	if u == nil {
		return false
	}

	switch v := u.(type) {
	case int:
		value := v
		if unit != value {
			return false
		}
	case interval:
		interval := u.(interval)
		if unit < interval.min || interval.max < unit {
			return false
		}
	}

	return true
}

func schemaValidation(sl validator.StructLevel) {
	schema := sl.Current().Interface().(entities.Schema)
	typeID := schema.TypeID

	if (typeID < 0 || 15 < typeID) && (typeID < 0xfff0 || 0xfff2 < typeID) && typeID != 0xff10 {
		sl.ReportError(schema, "schema", "Type ID", "typeID", "false")
		return
	}

	if !isValidValueType(schema.TypeID, schema.ValueType) {
		sl.ReportError(schema, "schema", "Value Type", "valueType", "false")
		return
	}

	if !isValidUnit(schema.TypeID, schema.Unit) {
		sl.ReportError(schema, "schema", "Unit", "unit", "false")
	}
}
