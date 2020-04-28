package interactors

import (
	"fmt"

	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
)

// PublishData executes the use case operations to publish data from the things to cloud
func (i *ThingInteractor) PublishData(authorization, thingID string, data []entities.Data) error {
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

	err = i.publisher.PublishPublishedData(thingID, data)
	if err != nil {
		return fmt.Errorf("error sending message to client: %w", err)
	}

	i.logger.Info("publish data message successfully sent")
	return nil
}
