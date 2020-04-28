package interactors

// Unregister runs the use case to remove a registered thing
func (i *ThingInteractor) Unregister(authorization, id string) error {
	i.logger.Debug("executing unregister thing use case")

	if authorization == "" {
		return ErrAuthNotProvided
	}

	if id == "" {
		return ErrIDNotProvided
	}

	err := i.thingProxy.Remove(authorization, id)
	if err != nil {
		sendErr := i.publisher.PublishUnregisteredDevice(id, err)
		if sendErr != nil {
			i.logger.Debug(err)
			return sendErr
		}
		return err
	}

	sendErr := i.publisher.PublishUnregisteredDevice(id, nil)
	if sendErr != nil {
		return sendErr
	}

	return nil
}
