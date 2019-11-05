package interactors

import "github.com/CESARBR/knot-babeltower/pkg/logging"

// CreateUser to interact to user
type CreateUser struct {
	logger logging.Logger
}

// NewCreateUser contructs the interactor
func NewCreateUser(logger logging.Logger) *CreateUser {
	return &CreateUser{logger}
}

// Execute runs the use case
func (cu *CreateUser) Execute() {
	// TODO: proxy message to user service
	cu.logger.Debug("Executing Create User interactor")
}
