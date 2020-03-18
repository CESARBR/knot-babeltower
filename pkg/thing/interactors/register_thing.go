package interactors

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	// ErrIDLength error occur when the thing's id have more than 16 ascii characters
	ErrIDLength = errors.New("id length exceeds 16 characters")

	// ErrIDNotInHex error occur when the thing's id is not formatted in hexadecimal base
	ErrIDNotInHex = errors.New("id is not in hexadecimal format")
)

// Register runs the use case to create a new thing
func (i *ThingInteractor) Register(authorization, id, name string) error {
	i.logger.Debug("Executing register thing use case")
	err := i.verifyThingID(id)
	if err != nil {
		sendErr := i.sendResponse(id, "", err)
		return fmt.Errorf("error registering thing: %w", sendErr)
	}

	// Get the id generated as a token and send in the response
	token, err := i.thingProxy.Create(id, name, authorization)
	sendErr := i.sendResponse(id, token, err)
	if err != nil {
		// it was not possible to create a thing, so returns without send request to connector
		return fmt.Errorf("error registering thing: %w", sendErr)
	}

	err = i.connectorPublisher.SendRegisterDevice(id, name)
	if err != nil {
		if sendErr != nil {
			// an error ocurred also on replying to client (on sendResponse method)
			return fmt.Errorf("error sending request to connector: %v: %w", err, sendErr)
		}
		return fmt.Errorf("error sending request to connector: %w", err)
	}

	return sendErr
}

func (i *ThingInteractor) verifyThingID(id string) error {
	if len(id) > 16 {
		return ErrIDLength
	}

	_, err := strconv.ParseUint(id, 16, 64)
	if err != nil {
		return ErrIDNotInHex
	}

	return nil
}

func (i *ThingInteractor) sendResponse(id, token string, err error) error {
	sendErr := i.clientPublisher.SendRegisteredDevice(id, token, getErrMessagePtr(err))
	if sendErr != nil {
		if err != nil {
			return fmt.Errorf("error sending response to client: %v: %w", sendErr, err)
		}
		return fmt.Errorf("error sending response to client: %w", sendErr)
	}
	return err
}

func getErrMessagePtr(err error) *string {
	if err != nil {
		msg := err.Error()
		return &msg
	}
	return nil
}
