package interactors

import (
	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/user/delivery/http"
	"github.com/CESARBR/knot-babeltower/pkg/user/entities"
)

// CreateToken has all dependencies and methods to enable the create token use case
// execution. It is composed by the logger and user proxy services.
type CreateToken struct {
	logger     logging.Logger
	usersProxy http.UsersProxy
	authProxy  http.AuthProxy
}

// NewCreateToken creates a new CreateToken instance by receiving its dependencies.
func NewCreateToken(logger logging.Logger, usersProxy http.UsersProxy, authProxy http.AuthProxy) *CreateToken {
	return &CreateToken{logger, usersProxy, authProxy}
}

// Execute receives the user entity filled with e-mail and password properties and try
// to create a token on the user proxy service. If it succeed, the token is returned.
func (ct *CreateToken) Execute(user entities.User, tokenType string, duration int) (token string, err error) {
	if tokenType == "user" {
		token, err = ct.usersProxy.CreateToken(user)
	} else if tokenType == "app" {
		token, err = ct.authProxy.CreateAppToken(user, duration)
	} else {
		err = entities.ErrInvalidTokenType
	}

	if err != nil {
		ct.logger.Errorf("failed to create a token: %s", err.Error())
		return "", err
	}

	return token, nil
}
