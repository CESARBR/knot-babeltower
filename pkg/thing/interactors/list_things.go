package interactors

import (
	"fmt"

	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
)

// List fetchs the registered things and return them as an array
func (i *ThingInteractor) List(authorization string) ([]*entities.Thing, error) {
	if authorization == "" {
		return nil, ErrAuthNotProvided
	}

	things, err := i.thingProxy.List(authorization)
	if err != nil {
		return nil, fmt.Errorf("error getting list of things: %w", err)
	}

	i.logger.Info("devices obtained")
	return things, nil
}
