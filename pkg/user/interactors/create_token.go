package interactors

import (
	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/user/delivery/http"
	"github.com/CESARBR/knot-babeltower/pkg/user/entities"
)

// CreateToken has all dependencies and methods to enable the create token use case
// execution. It is composed by the logger and user proxy services.
type CreateToken struct {
	logger    logging.Logger
	userProxy http.UserProxy
}

// NewCreateToken creates a new CreateToken instance by receiving its dependencies.
func NewCreateToken(logger logging.Logger, userProxy http.UserProxy) *CreateToken {
	return &CreateToken{logger, userProxy}
}

// Execute receives the user entity filled with e-mail and password properties and try
// to create a token on the user proxy service. If it succeed, the token is returned.
func (ct *CreateToken) Execute(user entities.User) (token string, err error) {
	token, err = ct.userProxy.CreateToken(user)
	if err != nil {
		ct.logger.Errorf("send request error: %s", err.Error())
		return "", err
	}

	return token, err
}
