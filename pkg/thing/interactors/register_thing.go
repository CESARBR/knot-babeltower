package interactors

import (
	"fmt"
	"strconv"

	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
)

const errorRegisteringThingMessage = "error registering thing"

// Register runs the use case to create a new thing
func (i *ThingInteractor) Register(authorization, id, name string) error {
	if id == "" {
		return ErrIDNotProvided
	}
	if name == "" {
		return ErrNameNotProvided
	}

	err := i.verifyThingID(id)
	if err != nil {
		sendErr := i.sendResponse(id, name, "", err)
		return fmt.Errorf("%s: %w", errorRegisteringThingMessage, sendErr)
	}

	// verify if thing is already registered
	_, err = i.thingProxy.Get(authorization, id)
	if err == nil {
		sendErr := i.sendResponse(id, name, "", entities.ErrThingExists)
		return fmt.Errorf("%s: %w", errorRegisteringThingMessage, sendErr)
	}

	// Get the id generated as a token and send in the response
	token, err := i.thingProxy.Create(id, name, authorization)
	sendErr := i.sendResponse(id, name, token, err)
	if err != nil {
		return fmt.Errorf("%s: %w", errorRegisteringThingMessage, sendErr)
	}

	return sendErr
}

func (i *ThingInteractor) verifyThingID(id string) error {
	if len(id) > 16 {
		return ErrIDLength
	}

	_, err := strconv.ParseUint(id, 16, 64)
	if err != nil {
		return ErrIDNotHex
	}

	return nil
}

func (i *ThingInteractor) sendResponse(id, name, token string, err error) error {
	sendErr := i.publisher.PublishRegisteredDevice(id, name, token, err)
	if sendErr != nil {
		if err != nil {
			return fmt.Errorf("error sending response to client: %v: %w", sendErr, err)
		}
		return fmt.Errorf("error sending response to client: %w", sendErr)
	}
	return err
}
