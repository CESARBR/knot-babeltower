package interactors

import (
	"fmt"

	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
)

var (
	// ErrNoAuthToken error occur when no authorization token is provided
	ErrNoAuthToken = fmt.Errorf("missing authorization token")

	// ErrNoIDParam error occur when no thing's id is provided
	ErrNoIDParam = fmt.Errorf("missing thing's id")

	// ErrNoDataParam error occur when no thing's data is provided
	ErrNoDataParam = fmt.Errorf("missing thing's data")

	// ErrNoSchema error occur when the thing has no schema yet
	ErrNoSchema = fmt.Errorf("thing has no schema")

	// ErrDataInvalid error occur when the provided data mismatch the thing's schema
	ErrDataInvalid = fmt.Errorf("data is incompatible with thing's schema")
)

// UpdateData executes the use case operations to update data in thing
func (i *ThingInteractor) UpdateData(authorization, thingID string, data []entities.Data) error {
	if authorization == "" {
		return ErrNoAuthToken
	}
	if thingID == "" {
		return ErrNoIDParam
	}
	if data == nil {
		return ErrNoDataParam
	}

	err := i.verifyThingData(authorization, thingID, data)
	if err != nil {
		return fmt.Errorf("error validating thing's data: %w", err)
	}

	err = i.clientPublisher.SendUpdateData(thingID, data)
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

	if thing.Schema == nil {
		return ErrNoSchema
	}

	for _, d := range data {
		if !validateSchema(d, thing.Schema) {
			return ErrDataInvalid
		}
	}

	return nil
}

func validateSchema(data entities.Data, schema []entities.Schema) bool {
	for _, s := range schema {
		if s.SensorID == data.SensorID {
			switch data.Value.(type) {
			case int:
				return s.ValueType == 1 // int
			case float64:
				return s.ValueType == 2 // float
			case bool:
				return s.ValueType == 3 // bool
			case string:
				return s.ValueType == 4 // raw
			default:
				return false
			}
		}
	}

	return false
}
