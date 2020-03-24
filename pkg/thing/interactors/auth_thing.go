package interactors

import "fmt"

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
		sendErr := i.clientPublisher.SendAuthStatus(id, err)
		if sendErr != nil {
			return fmt.Errorf("error getting thing metadata: %v: %w", sendErr, err)
		}
		return fmt.Errorf("error getting thing metadata: %w", err)
	}

	err = i.clientPublisher.SendAuthStatus(id, nil)
	if err != nil {
		return fmt.Errorf("error sending response to client: %w", err)
	}

	i.logger.Info("authentication status sucessfully sent")
	return nil
}
