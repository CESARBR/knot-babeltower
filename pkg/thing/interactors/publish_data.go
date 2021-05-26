package interactors

import (
	"fmt"

	"github.com/CESARBR/knot-babeltower/pkg/jwt"
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

	err = i.publisher.PublishBroadcastData(thingID, authorization, data)
	if err != nil {
		return fmt.Errorf("error publishing data in broadcast mode: %w", err)
	}

	err = i.publishSessionData(thingID, authorization, data)
	if err != nil {
		return fmt.Errorf("error publishing data to user sessions: %w", err)
	}

	return nil
}

func (i *ThingInteractor) publishSessionData(thingID, authorization string, data []entities.Data) error {
	email, err := jwt.GetEmail(authorization)
	if err != nil {
		return fmt.Errorf("error getting user e-mail from token: %w", err)
	}

	sessionId, err := i.sessionStore.Get(email)
	if err != nil {
		return fmt.Errorf("error getting user session: %w", err)
	}

	if sessionId != "" {
		err = i.publisher.PublishSessionData(thingID, authorization, sessionId, data)
		if err != nil {
			return fmt.Errorf("error sending message to client: %w", err)
		}
	}

	return nil
}
