package interactors

// Unregister runs the use case to remove a registered thing
func (i *ThingInteractor) Unregister(authorization, id string) error {
	i.logger.Debug("executing unregister thing use case")

	if id == "" {
		return ErrIDNotProvided
	}

	err := i.thingProxy.Remove(authorization, id)
	if err != nil {
		sendErr := i.publisher.PublishUnregisteredDevice(id, authorization, err)
		if sendErr != nil {
			i.logger.Debug(err)
			return sendErr
		}
		return err
	}

	sendErr := i.publisher.PublishUnregisteredDevice(id, authorization, nil)
	if sendErr != nil {
		return sendErr
	}

	return nil
}
