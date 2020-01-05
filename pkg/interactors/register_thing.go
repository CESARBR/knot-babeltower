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
}

// ErrorIDLenght is raised when ID is more than 16 characters
type ErrorIDLenght struct{}

// ErrorIDInvalid is raised when ID is not in hexadecimal value
type ErrorIDInvalid struct{}

// ErrorNameNotFound is raised when Name is empty
type ErrorNameNotFound struct{}

// ErrorArgument is raised when Name is empty
type ErrorArgument struct{ msg string }

func (err ErrorIDLenght) Error() string {
	return "ID length error"
}

func (err ErrorIDInvalid) Error() string {
	return "ID is not in hexadecimal"
}

func (err ErrorNameNotFound) Error() string {
	return "Name not found"
}

func (err ErrorArgument) Error() string {
	return err.msg
}

// NewRegisterThing contructs the use case
func NewRegisterThing(logger logging.Logger, msgPublisher network.Publisher) *RegisterThing {
	return &RegisterThing{logger, msgPublisher}
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

func (rt *RegisterThing) verifyArguments(args ...interface{}) error {
	if len(args) < 1 {
		return ErrorArgument{"Missing argument name"}
	}

	name, ok := args[0].(string)
	if !ok {
		return ErrorArgument{msg: "Name is not string"}
	}

	if len(name) == 0 {
		return ErrorNameNotFound{}
	}

	return nil
}

// Execute runs the use case
func (rt *RegisterThing) Execute(id string, args ...interface{}) error {
	rt.logger.Debug("Executing register thing use case")
	err := rt.verifyArguments(args...)
	if err != nil {
		errReply := rt.reply(id, "", err)
		if errReply != nil {
			rt.logger.Error(errReply)
			return errReply
		}

		return err
	}

	err = rt.verifyThingID(id)
	errReply := rt.reply(id, "", err)
	if errReply != nil {
		rt.logger.Error(errReply)
		return errReply
	}

	// TODO: add proxy request to token

	return nil
}
