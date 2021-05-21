package interactors

import (
	"fmt"

	"github.com/CESARBR/knot-babeltower/pkg/cache"
	"github.com/CESARBR/knot-babeltower/pkg/jwt"
	"github.com/CESARBR/knot-babeltower/pkg/thing/delivery/http"
	"github.com/CESARBR/knot-babeltower/pkg/user/entities"
)

// CreateSession is a use case operation that creates a session for the user.
// This session is represented by a random ID, which will be further used for
// publishing thing's data.
type CreateSession struct {
	thingsProxy  http.ThingProxy
	generator    Generator
	sessionStore cache.SessionStore
}

// NewCreateSession creates a new CreateSession instance by receiving its dependencies.
func NewCreateSession(
	thingsProxy http.ThingProxy,
	generator Generator,
	sessionStore cache.SessionStore,
) *CreateSession {
	return &CreateSession{thingsProxy, generator, sessionStore}
}

// Execute creates a new user session by receiving the authorization token as parameter.
// This function will verify if a session already exists for the users. If not, it will
// create a new session ID and save it in some caching storage.
func (cs *CreateSession) Execute(authorization string) (string, error) {
	email, err := jwt.GetEmail(authorization)
	if err != nil {
		return "", fmt.Errorf("failed to get user email from authorization token: %w", err)
	}

	if !cs.isValidToken(authorization) {
		return "", entities.ErrTokenForbidden
	}

	session, err := cs.sessionStore.Get(email)
	if err != nil {
		return "", fmt.Errorf("failed to get user session: %w", err)
	}

	if session != "" {
		return session, nil
	}

	id := cs.generator.ID()
	err = cs.sessionStore.Save(email, id)
	if err != nil {
		return "", fmt.Errorf("failed to save user session: %w", err)
	}

	return id, err
}

func (cs *CreateSession) isValidToken(authorization string) bool {
	_, err := cs.thingsProxy.List(authorization)
	return err == nil
}
