package interactors

import (
	"fmt"
	"math"

	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
)

// UpdateConfig executes the use case to update thing's configuration
func (i *ThingInteractor) UpdateConfig(authorization, id string, configList []entities.Config) error {
	if authorization == "" {
		return ErrAuthNotProvided
	}
	if id == "" {
		return ErrIDNotProvided
	}
	if configList == nil {
		return ErrConfigNotProvided
	}

	err := i.validateConfig(authorization, id, configList)
	if err != nil {
		return fmt.Errorf("failed to validate if config is valid: %w", err)
	}

	err = i.thingProxy.UpdateConfig(authorization, id, configList)
	if err != nil {
		return err
	}

	return nil
}

func (i *ThingInteractor) validateConfig(authorization, id string, configList []entities.Config) error {
	thing, err := i.thingProxy.Get(authorization, id)
	if err != nil {
		return fmt.Errorf("error getting thing metadata: %w", err)
	}

	if thing.Schema == nil {
		return ErrSchemaUndefined
	}

	err = validateSchemaMatch(configList, thing.Schema)
	if err != nil {
		return err
	}

	return nil
}

func validateSchemaMatch(configList []entities.Config, schemaList []entities.Schema) error {
	for _, c := range configList {
		schema := checkSensorInSchema(c.SensorID, schemaList)
		if schema == nil {
			return ErrConfigInvalid
		}

		err := validateFlagValue(c, schema.ValueType)
		if err != nil {
			return err
		}
	}

	return nil
}

func checkSensorInSchema(sensorID int, schemaList []entities.Schema) *entities.Schema {
	for _, s := range schemaList {
		if sensorID == s.SensorID {
			return &s
		}
	}

	return nil
}

func validateFlagValue(config entities.Config, valueType int) error {
	if config.LowerThreshold != nil && !isValidValue(config.LowerThreshold, valueType) {
		return ErrDataInvalid
	}

	if config.UpperThreshold != nil && !isValidValue(config.UpperThreshold, valueType) {
		return ErrDataInvalid
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
