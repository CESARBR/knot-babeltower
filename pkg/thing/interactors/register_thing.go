package interactors

import (
	"strconv"

	"github.com/CESARBR/knot-babeltower/pkg/network"
)

// ErrorIDLenght is raised when ID is more than 16 characters
type ErrorIDLenght struct{}

// ErrorIDInvalid is raised when ID is not in hexadecimal value
type ErrorIDInvalid struct{}

// ErrorNameNotFound is raised when Name is empty
type ErrorNameNotFound struct{}

// ErrorMissingArgument is raised there is some argument missing
type ErrorMissingArgument struct{}

// ErrorInvalidTypeArgument is raised when the type is the expected
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

func (i *ThingInteractor) reply(id, token string, err error) error {
	var errStr *string

	if err != nil {
		errStr = new(string)
		*errStr = err.Error()
	} else {
		errStr = nil
	}

	response := network.RegisterResponseMsg{ID: id, Token: token, Error: errStr}
	errPublish := i.clientPublisher.SendRegisterDevice(response)
	if errPublish != nil {
		i.logger.Error(errPublish)
		return errPublish
	}

	return nil
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

// Register runs the use case to create a new thing
func (i *ThingInteractor) Register(authorization, id, name string) error {
	i.logger.Debug("Executing register thing use case")
	err := i.verifyThingID(id)
	if err != nil {
		errReply := i.reply(id, "", err)
		if errReply != nil {
			i.logger.Error(errReply)
		}

		return err
	}

	// Get the id generated as a token and send in the response
	token, err := i.thingProxy.Create(id, name, authorization)
	errReply := i.reply(id, token, err)
	if errReply != nil {
		i.logger.Error(errReply)
		return errReply
	}

	err = i.connectorPublisher.SendRegisterDevice(id, name)
	if err != nil {
		i.logger.Error(err)
		return err
	}

	return nil
}
