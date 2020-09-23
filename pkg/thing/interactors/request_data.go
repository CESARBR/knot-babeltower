package interactors

import (
	"fmt"

	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
)

// RequestData executes the use case operations to request data from the thing
func (i *ThingInteractor) RequestData(authorization, thingID string, sensorIds []int) error {
	if authorization == "" {
		return ErrAuthNotProvided
	}
	if thingID == "" {
		return ErrIDNotProvided
	}
	if sensorIds == nil {
		return ErrSensorsNotProvided
	}

	thing, err := i.thingProxy.Get(authorization, thingID)
	if err != nil {
		i.logger.Error(err)
		return err
	}

	if thing.Config == nil {
		i.logger.Error(fmt.Errorf("thing %s has no config yet", thing.ID))
		return err
	}

	err = validateSensors(sensorIds, thing.Config)
	if err != nil {
		i.logger.Error(err)
		return err
	}

	err = i.publisher.PublishRequestData(thingID, sensorIds)
	if err != nil {
		i.logger.Error(err)
		return err
	}

	i.logger.Info("data request command successfully sent")
	return nil
}

// validateSensors validates a slice of sensor ids against the thing's registered schema
// that represents the sensors and actuators associated to it.
func validateSensors(sensorIds []int, configList []entities.Config) error {
	for _, id := range sensorIds {
		if !sensorExists(configList, id) {
			return ErrSensorInvalid
		}
	}

	return nil
}

func sensorExists(configList []entities.Config, id int) bool {
	for _, c := range configList {
		if c.SensorID == id {
			return true
		}
	}

	return false
}
