package interactors

import (
	"errors"
)

// List fetchs the registered things and return them as an array
func (i *ThingInteractor) List(authorization string) error {
	if authorization == "" {
		return errors.New("authorization key not provided")
	}

	things, err := i.thingProxy.List(authorization)
	if err != nil {
		return err
	}

	err = i.clientPublisher.SendDevicesList(things)
	if err != nil {
		return err
	}

	i.logger.Info("devices obtained")
	return nil
}
