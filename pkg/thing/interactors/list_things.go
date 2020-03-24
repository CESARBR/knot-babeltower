package interactors

import "fmt"

// List fetchs the registered things and return them as an array
func (i *ThingInteractor) List(authorization string) error {
	if authorization == "" {
		return ErrAuthNotProvided
	}

	things, err := i.thingProxy.List(authorization)
	if err != nil {
		return fmt.Errorf("error getting list of things: %w", err)
	}

	err = i.clientPublisher.SendDevicesList(things)
	if err != nil {
		return fmt.Errorf("error sending response to client: %w", err)
	}

	i.logger.Info("devices obtained")
	return nil
}
