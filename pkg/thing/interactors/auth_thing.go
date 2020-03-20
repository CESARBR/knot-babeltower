package interactors

// Auth is responsible to implement the thing's authentication use case
func (i *ThingInteractor) Auth(authorization, id string) error {
	if authorization == "" {
		return ErrAuthNotProvided
	}
	if id == "" {
		return ErrIDNotProvided
	}

	_, err := i.thingProxy.Get(authorization, id)
	if err != nil {
		i.logger.Error(err)
		err = i.clientPublisher.SendAuthStatus(id, err)
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
