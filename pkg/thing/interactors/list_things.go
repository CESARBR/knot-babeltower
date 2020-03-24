package interactors

import "fmt"

// List fetchs the registered things and return them as an array
func (i *ThingInteractor) List(authorization string) error {
	if authorization == "" {
		return ErrAuthNotProvided
	}

	things, err := i.thingProxy.List(authorization)
	sendErr := i.clientPublisher.SendDevicesList(things, err)
	if err != nil {
		if sendErr != nil {
			return fmt.Errorf("error getting list of things: %v, %w", sendErr, err)
		}
		return fmt.Errorf("error getting list of things: %w", err)
	}
	if sendErr != nil {
		return fmt.Errorf("error sending response to client: %w", sendErr)
	}

	i.logger.Info("devices obtained")
	return nil
}
