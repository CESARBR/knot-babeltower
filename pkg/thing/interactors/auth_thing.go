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
		return fmt.Errorf("can't receive thing metadata: %w", err)
	}

	return nil
}
