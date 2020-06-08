package interactors

import (
	"errors"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/user/delivery/http"
	"github.com/CESARBR/knot-babeltower/pkg/user/entities"
)

// CreateToken has all dependencies and methods to enable the create token use case
// execution. It is composed by the logger and user proxy services.
type CreateToken struct {
	logger     logging.Logger
	usersProxy http.UsersProxy
	authnProxy http.AuthnProxy
}

// NewCreateToken creates a new CreateToken instance by receiving its dependencies.
func NewCreateToken(logger logging.Logger, usersProxy http.UsersProxy, authnProxy http.AuthnProxy) *CreateToken {
	return &CreateToken{logger, usersProxy, authnProxy}
}

// Execute receives the user entity filled with e-mail and password properties and try
// to create a token on the user proxy service. If it succeed, the token is returned.
func (ct *CreateToken) Execute(user entities.User, tokenType string) (token string, err error) {
	if tokenType == "user" {
		token, err = ct.usersProxy.CreateToken(user)
	} else if tokenType == "app" {
		token, err = ct.authnProxy.CreateAppToken(user)
	} else {
		err = errors.New("only 'user' and 'app' token types are supported")
	}

	if err != nil {
		ct.logger.Errorf("failed to create a token: %s", err.Error())
		return "", err
	}

	return token, nil
}
