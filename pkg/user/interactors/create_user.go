package interactors

import (
	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/user/delivery/http"
	"github.com/CESARBR/knot-babeltower/pkg/user/entities"
)

// CreateUser to interact to user
type CreateUser struct {
	logger    logging.Logger
	userProxy http.UserProxy
}

// NewCreateUser contructs the interactor
func NewCreateUser(logger logging.Logger, userProxy http.UserProxy) *CreateUser {
	return &CreateUser{logger, userProxy}
}

// Execute runs the use case
func (cu *CreateUser) Execute(user entities.User) (err error) {
	cu.logger.Debug("executing Create User interactor")

	err = cu.userProxy.Create(user)
	if err != nil {
		cu.logger.Errorf("send request error: %s", err.Error())
	}

	return err
}
