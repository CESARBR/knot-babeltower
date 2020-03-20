package interactors

import (
	"errors"
)

// Unregister runs the use case to remove a registered thing
func (i *ThingInteractor) Unregister(authorization, id string) error {
	i.logger.Debug("Executing unregister thing use case")

	if authorization == "" {
		return errors.New("authorization key not provided")
	}

	if id == "" {
		return errors.New("thing's id not provided")
	}

	err := i.thingProxy.Remove(authorization, id)
	if err != nil {
		sendErr := i.clientPublisher.SendUnregisteredDevice(id, err)
		if sendErr != nil {
			i.logger.Debug(err)
			return sendErr
		}
		return err
	}

	sendErr := i.connectorPublisher.SendUnregisterDevice(id)
	if sendErr != nil {
		// TODO handle error at sending unregister message to connector
		i.logger.Error(sendErr)
	}

	sendErr = i.clientPublisher.SendUnregisteredDevice(id, nil)
	if sendErr != nil {
		return sendErr
	}

	return nil
}
