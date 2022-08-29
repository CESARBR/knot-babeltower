package interactors

import (
	"fmt"
	"math"

	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
)

const (
	intType    = 1
	floatType  = 2
	boolType   = 3
	rawType    = 4
	int64Type  = 5
	uintType   = 6
	uint64Type = 7
)

// UpdateData executes the use case operations to update data in thing
func (i *ThingInteractor) UpdateData(authorization, thingID string, data []entities.Data) error {
	if authorization == "" {
		return ErrAuthNotProvided
	}
	if thingID == "" {
		return ErrIDNotProvided
	}
	if data == nil {
		return ErrDataNotProvided
	}

	err := i.verifyThingData(authorization, thingID, data)
	if err != nil {
		return fmt.Errorf("error validating thing's data: %w", err)
	}

	err = i.publisher.PublishUpdateData(thingID, data)
	if err != nil {
		return fmt.Errorf("error sending message to client: %w", err)
	}

	i.logger.Info("data update command successfully sent")
	return nil
}

func (i *ThingInteractor) verifyThingData(authorization, thingID string, data []entities.Data) error {
	thing, err := i.thingProxy.Get(authorization, thingID)
	if err != nil {
		return fmt.Errorf("error getting thing metadata: %w", err)
	}

	if thing.Config == nil {
		return ErrConfigUndefined
	}

	for _, d := range data {
		if !validateSchema(d, thing.Config) {
			return ErrDataInvalid
		}
	}

	return nil
}

func validateSchema(data entities.Data, configList []entities.Config) bool {
	for _, c := range configList {
		if c.SensorID == data.SensorID {
			switch data.Value.(type) {
			case float64:
				if data.Value == math.Trunc(data.Value.(float64)) { // check if number is integer
					return ValidateSchemaNumber(data.Value.(float64), c.Schema.ValueType)
				}
				return c.Schema.ValueType == floatType
			case bool:
				return c.Schema.ValueType == boolType
			case string:
				return c.Schema.ValueType == rawType
			default:
				return false
			}
		}
	}

	return false
}

func createValueValidatorMapping() map[int]func(value float64) bool {
	valueValidatorMapping := make(map[int]func(value float64) bool)
	valueValidatorMapping[intType] = isValidInt
	valueValidatorMapping[floatType] = isValidFloat
	valueValidatorMapping[int64Type] = isValidInt64
	valueValidatorMapping[uintType] = isValidUint
	valueValidatorMapping[uint64Type] = isValidUint64
	return valueValidatorMapping
}

func isValidInt(value float64) bool {
	return value >= math.MinInt32 && value <= math.MaxInt32
}
func isValidFloat(value float64) bool {
	return true
}
func isValidInt64(value float64) bool {
	return value >= math.MinInt64 && value <= math.MaxInt64
}
func isValidUint(value float64) bool {
	return value >= 0 && value <= math.MaxUint32
}
func isValidUint64(value float64) bool {
	return value >= 0 && value <= math.MaxUint64
}

// Validates the value received against its type defined in the sensor's schema
func ValidateSchemaNumber(value float64, valueType int) bool {
	valueValidatorMapping := createValueValidatorMapping()
	if function, ok := valueValidatorMapping[valueType]; ok {
		return function(value)
	} else {
		return false
	}
}
