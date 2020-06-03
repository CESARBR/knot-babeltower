package interactors

import (
	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/user/delivery/http"
	"github.com/CESARBR/knot-babeltower/pkg/user/entities"
)

// CreateUser to interact to user
type CreateUser struct {
	logger     logging.Logger
	usersProxy http.UsersProxy
}

// NewCreateUser contructs the interactor
func NewCreateUser(logger logging.Logger, usersProxy http.UsersProxy) *CreateUser {
	return &CreateUser{logger, usersProxy}
}

// Execute runs the use case
func (cu *CreateUser) Execute(user entities.User) (err error) {
	err = cu.usersProxy.Create(user)
	if err != nil {
		cu.logger.Errorf("failed to create a new user: %s", err.Error())
	}

	return err
}
