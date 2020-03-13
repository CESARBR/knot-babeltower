package interactors

import (
	"fmt"
	"strconv"
)

// ErrorIDLenght is raised when ID is more than 16 characters
type ErrorIDLenght struct{}

// ErrorIDInvalid is raised when ID is not in hexadecimal value
type ErrorIDInvalid struct{}

// ErrorNameNotFound is raised when Name is empty
type ErrorNameNotFound struct{}

// ErrorMissingArgument is raised when there is some argument missing
type ErrorMissingArgument struct{}

// ErrorInvalidTypeArgument is raised when the type is not the expected
type ErrorInvalidTypeArgument struct{ msg string }

func (err ErrorIDLenght) Error() string {
	return "ID length error"
}

func (err ErrorIDInvalid) Error() string {
	return "ID is not in hexadecimal"
}

func (err ErrorNameNotFound) Error() string {
	return "Name not found"
}

func (err ErrorMissingArgument) Error() string {
	return "Missing arguments"
}

func (err ErrorInvalidTypeArgument) Error() string {
	return err.msg
}

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
		return ErrorIDLenght{}
	}

	_, err := strconv.ParseUint(id, 16, 64)
	if err != nil {
		i.logger.Error(err)
		return ErrorIDInvalid{}
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
