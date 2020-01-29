package interactors

import (
	"errors"
)

// Auth is responsible to implement the thing's authentication use case
func (i *ThingInteractor) Auth(authorization, id, token string) error {

	if authorization == "" {
		return errors.New("authorization key not provided")
	}

	if id == "" {
		return errors.New("thing's id not provided")
	}

	if token == "" {
		return errors.New("thing's token not provided")
	}

	_, err := i.thingProxy.GetThing(authorization, id)
	if err != nil {
		i.logger.Error(err)
		msg := err.Error()
		err = i.clientPublisher.SendAuthStatus(id, &msg)
		return err
	}

	err = i.clientPublisher.SendAuthStatus(id, nil)
	if err != nil {
		i.logger.Error(err)
		return err
	}

	i.logger.Info("authentication status sucessfully sent")
	return nil
}
