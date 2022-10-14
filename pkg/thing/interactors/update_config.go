package interactors

import (
	"fmt"
	"math"
	"reflect"

	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
	"github.com/go-playground/validator"
)

type schemaType struct {
	valueType interface{}
	unit      interface{}
}

type interval struct {
	min int
	max int
}

// rules reference table: https://knot-devel.cesar.org.br/doc/thing/unit-type-value.html
var rules = map[int]schemaType{
	0x0000: {valueType: interval{1, 8}, unit: 0},              // NONE
	0x0001: {valueType: interval{1, 8}, unit: interval{1, 3}}, // VOLTAGE
	0x0002: {valueType: interval{1, 8}, unit: interval{1, 2}}, // CURRENT
	0x0003: {valueType: interval{1, 8}, unit: 1},              // RESISTENCE
	0x0004: {valueType: interval{1, 8}, unit: interval{1, 3}}, // POWER
	0x0005: {valueType: interval{1, 8}, unit: interval{1, 3}}, // TEMPERATURE
	0x0006: {valueType: interval{1, 8}, unit: 1},              // RELATIVE_HUMIDITY
	0x0007: {valueType: interval{1, 8}, unit: interval{1, 3}}, // LUMINOSITY
	0x0008: {valueType: interval{1, 8}, unit: interval{1, 3}}, // TIME
	0x0009: {valueType: interval{1, 8}, unit: interval{1, 4}}, // MASS
	0x000A: {valueType: interval{1, 8}, unit: interval{1, 3}}, // PRESSURE
	0x000B: {valueType: interval{1, 8}, unit: interval{1, 4}}, // DISTANCE
	0x000C: {valueType: interval{1, 8}, unit: interval{1, 2}}, // ANGLE
	0x000D: {valueType: interval{1, 8}, unit: interval{1, 4}}, // VOLUME
	0x000E: {valueType: interval{1, 8}, unit: interval{1, 3}}, // AREA
	0x000F: {valueType: interval{1, 8}, unit: 1},              // RAIN
	0x0010: {valueType: interval{1, 8}, unit: 1},              // DENSITY
	0x0011: {valueType: interval{1, 8}, unit: 1},              // LATITUDE
	0x0012: {valueType: interval{1, 8}, unit: 1},              // LONGITUDE
	0x0013: {valueType: interval{1, 8}, unit: interval{1, 4}}, // SPEED
	0x0014: {valueType: interval{1, 8}, unit: interval{1, 6}}, // VOLUMEFLOW
	0x0015: {valueType: interval{1, 8}, unit: interval{1, 6}}, // ENERGY
	0xFFF0: {valueType: interval{1, 8}, unit: 0},              // PRESENCE
	0xFFF1: {valueType: interval{1, 8}, unit: 0},              // SWITCH
	0xFFF2: {valueType: interval{1, 8}, unit: 0},              // COMMAND
	0xFF10: {valueType: interval{1, 8}, unit: interval{0, 1}}, // GENERIC
	0xFFFF: {valueType: interval{1, 8}, unit: 0},              // INVALID
}

// UpdateConfig executes the use case to update thing's configuration
// It returns two values:
//   - error: indicates if something goes wrong
//   - bool: indicates if the operation has changed something in the current thing's configuration
func (i *ThingInteractor) UpdateConfig(authorization, id string, configList []entities.Config) (bool, error) {
	if id == "" {
		return false, ErrIDNotProvided
	}
	if configList == nil {
		return false, ErrConfigNotProvided
	}

	err := i.validateConfig(authorization, id, configList)
	if err != nil {
		if err == ErrConfigEqual {
			return false, nil
		}

		return false, fmt.Errorf("failed to validate if config is valid: %w", err)
	}

	err = i.thingProxy.UpdateConfig(authorization, id, configList)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (i *ThingInteractor) validateConfig(authorization, id string, configList []entities.Config) error {
	thing, err := i.thingProxy.Get(authorization, id)
	if err != nil {
		return fmt.Errorf("error getting thing metadata: %w", err)
	}

	err = validateSchemaExists(configList, thing.Config)
	if err != nil {
		return err
	}

	if reflect.DeepEqual(thing.Config, configList) {
		return ErrConfigEqual
	}

	if !i.isValidSchema(configList) {
		return ErrSchemaInvalid
	}
	configList = validateConfigIntegrity(configList, thing.Config)

	err = validateFlagValue(configList)
	if err != nil {
		return err
	}

	return nil
}

func validateFlagValue(configList []entities.Config) error {
	for _, c := range configList {
		if c.Event.LowerThreshold != nil && !isValidValue(c.Event.LowerThreshold, c.Schema.ValueType) {
			return ErrDataInvalid
		}

		if c.Event.UpperThreshold != nil && !isValidValue(c.Event.UpperThreshold, c.Schema.ValueType) {
			return ErrDataInvalid
		}
	}

	return nil
}

func isValidValue(value interface{}, valueType int) bool {
	switch value := value.(type) {
	case float64:
		if value == math.Trunc(value) {
			return ValidateSchemaNumber(value, valueType)
		}
		return valueType == 2
	default:
		return false
	}
}

func validateSchemaExists(newConfigList []entities.Config, actualConfigList []entities.Config) error {
	for _, c := range newConfigList {
		if isAnExistentConfig(actualConfigList, c) {
			return ErrSchemaNotProvided
		}
	}

	return nil
}

func isAnExistentConfig(configList []entities.Config, config entities.Config) bool {
	hasRegisteredConfig := false
	for _, c := range configList {
		if config.SensorID == c.SensorID {
			hasRegisteredConfig = true
		}
	}

	if !hasRegisteredConfig && isSchemaEmpty(config.Schema) {
		return true
	}
	return false
}

func validateConfigIntegrity(newConfigList []entities.Config, actualConfigList []entities.Config) []entities.Config {
	for index, c := range newConfigList {
		for _, t := range actualConfigList {
			var event *entities.Event = &c.Event
			if c.SensorID == t.SensorID {
				if isSchemaEmpty(c.Schema) {
					newConfigList[index].Schema = t.Schema
				}
				if isEventEmpty(c.Event) {
					newConfigList[index].Event = t.Event
				}
				if c.Schema.ValueType != t.Schema.ValueType && event != nil && t.Event.LowerThreshold != nil && t.Event.UpperThreshold != nil {
					newConfigList[index].Event.LowerThreshold = nil
					newConfigList[index].Event.UpperThreshold = nil
				}
			}
		}
	}

	return newConfigList
}

func (i *ThingInteractor) isValidSchema(configList []entities.Config) bool {
	validate := validator.New()
	validate.RegisterStructValidation(schemaValidation, entities.Schema{})
	for _, config := range configList {
		if !isSchemaEmpty(config.Schema) {
			err := validate.Struct(config.Schema)
			if err != nil {
				return false
			}
		}
	}

	return true
}

func schemaValidation(sl validator.StructLevel) {
	fmt.Print("Running schemaValidation")
	schema := sl.Current().Interface().(entities.Schema)
	typeID := schema.TypeID

	if (typeID < 0 || 15 < typeID) && (typeID < 0xfff0 || 0xfff2 < typeID) && typeID != 0xff10 {
		sl.ReportError(schema, "schema", "Type ID", "typeID", "false")
		return
	}
	fmt.Print("passed typeID")
	if !isValidValueType(schema.TypeID, schema.ValueType) {
		sl.ReportError(schema, "schema", "Value Type", "valueType", "false")
		return
	}
	fmt.Print("passed valueType")
	if !isValidUnit(schema.TypeID, schema.Unit) {
		sl.ReportError(schema, "schema", "Unit", "unit", "false")
	}
	fmt.Print("passed unit")
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

func isSchemaEmpty(schema entities.Schema) bool {
	if schema.Name == "" && schema.TypeID == 0 && schema.Unit == 0 && schema.ValueType == 0 {
		return true
	}
	return false
}

func isEventEmpty(event entities.Event) bool {
	if !event.Change && event.TimeSec == 0 && event.LowerThreshold == nil && event.UpperThreshold == nil {
		return true
	}
	return false
}
