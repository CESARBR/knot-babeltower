package interactors

import (
	"fmt"

	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
)

// PublishData executes the use case operations to publish data from the things to cloud
func (i *ThingInteractor) PublishData(authorization, thingID string, data []entities.Data) error {
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

	err = i.connectorPublisher.SendPublishData(thingID, data)
	if err != nil {
		return fmt.Errorf("error sending message to connector: %w", err)
	}

	i.logger.Info("publish data message successfully sent")
	return nil
}
