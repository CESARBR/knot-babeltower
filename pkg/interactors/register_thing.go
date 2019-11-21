package interactors

import (
	"strconv"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/network"
)

// RegisterThing use case to register a new thing
type RegisterThing struct {
	logger       logging.Logger
	msgPublisher network.Publisher
	thingProxy   network.ThingProxy
}

// ErrorIDLenght is raised when ID is more than 16 characters
type ErrorIDLenght struct{}

// ErrorIDInvalid is raised when ID is not in hexadecimal value
type ErrorIDInvalid struct{}

// ErrorNameNotFound is raised when Name is empty
type ErrorNameNotFound struct{}

// ErrorUnauthorized is raised when authorization token is empty
type ErrorUnauthorized struct{}

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

func (err ErrorUnauthorized) Error() string {
	return "Authorization token not found"
}

func (err ErrorMissingArgument) Error() string {
	return "Missing arguments"
}

func (err ErrorInvalidTypeArgument) Error() string {
	return err.msg
}

// NewRegisterThing contructs the use case
func NewRegisterThing(logger logging.Logger, msgPublisher network.Publisher, thingProxy network.ThingProxy) *RegisterThing {
	return &RegisterThing{logger, msgPublisher, thingProxy}
}

func (rt *RegisterThing) reply(id, token string, err error) error {
	var errStr *string

	if err != nil {
		errStr = new(string)
		*errStr = err.Error()
	} else {
		errStr = nil
	}

	response := network.RegisterResponseMsg{ID: id, Token: token, Error: errStr}
	errPublish := rt.msgPublisher.SendRegisterDevice(response)
	if errPublish != nil {
		rt.logger.Error(errPublish)
		return errPublish
	}

	return nil
}

func (rt *RegisterThing) verifyThingID(id string) error {
	if len(id) > 16 {
		return ErrorIDLenght{}
	}

	_, err := strconv.ParseUint(id, 16, 64)
	if err != nil {
		rt.logger.Error(err)
		return ErrorIDInvalid{}
	}

	return nil
}

func (rt *RegisterThing) getArguments(args ...interface{}) (string, string, error) {
	if len(args) < 2 {
		return "", "", ErrorMissingArgument{}
	}

	name, ok := args[0].(string)
	if !ok {
		return "", "", ErrorInvalidTypeArgument{msg: "Name is not string"}
	}

	if len(name) == 0 {
		return "", "", ErrorNameNotFound{}
	}

	authorizationToken, ok := args[1].(string)
	if !ok {
		return "", "", ErrorInvalidTypeArgument{msg: "Authorization token is not string"}
	}

	if len(authorizationToken) == 0 {
		return "", "", ErrorUnauthorized{}
	}

	return name, authorizationToken, nil
}

// Execute runs the use case
func (rt *RegisterThing) Execute(id string, args ...interface{}) error {
	rt.logger.Debug("Executing register thing use case")
	name, authorizationToken, err := rt.getArguments(args...)
	if err != nil {
		errReply := rt.reply(id, "", err)
		if errReply != nil {
			rt.logger.Error(errReply)
			return errReply
		}

		return err
	}

	err = rt.verifyThingID(id)
	if err != nil {
		errReply := rt.reply(id, "", err)
		if errReply != nil {
			rt.logger.Error(errReply)
		}

		return err
	}

	// Get the id generated as a token and send in the response
	token, err := rt.thingProxy.Create(id, name, authorizationToken)
	errReply := rt.reply(id, token, err)
	if errReply != nil {
		rt.logger.Error(errReply)
		return errReply
	}

	return nil
}
