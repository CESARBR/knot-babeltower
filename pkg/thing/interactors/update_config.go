package interactors

import (
	"fmt"

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

	for _, c := range configList {
		schema := checkSensorInSchema(c.SensorID, thing.Schema)
		if schema == nil {
			return ErrConfigInvalid
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
