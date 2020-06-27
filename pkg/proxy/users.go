package proxy

import (
	"fmt"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/network"
	"github.com/CESARBR/knot-babeltower/pkg/user/entities"
)

// users documentation: https://github.com/mainflux/mainflux/blob/master/users/swagger.yaml

// UsersProxy represents the interface to the user's proxy operations
type UsersProxy interface {
	Create(user entities.User) (err error)
	CreateToken(user entities.User) (token string, err error)
}

// Users is responsible for implementing the user's proxy operations
type Users struct {
	URL    string
	http   *network.HTTP
	logger logging.Logger
}

// NewUsersProxy creates a new Proxy instance
func NewUsersProxy(logger logging.Logger, http *network.HTTP, userHost string, userPort uint16) UsersProxy {
	URL := fmt.Sprintf("http://%s:%d", userHost, userPort)
	logger.Debug("user proxy configured to " + URL)
	return &Users{URL, http, logger}
}

// userSchema represents the schema for an user
type userSchema struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// tokenSchema represents the schema for a token
type tokenSchema struct {
	Token string `json:"token"`
}

// Create create a new user on users service
func (u *Users) Create(user entities.User) error {
	request := network.Request{
		Path:   u.URL + "/users",
		Method: "POST",
		Body:   userSchema{Email: user.Email, Password: user.Password},
	}

	err := u.http.MakeRequest(request, nil, StatusErrors)
	if err != nil {
		return fmt.Errorf("error creating a new user: %w", err)
	}

	return nil
}

// CreateToken creates a valid token for the specified user
func (u *Users) CreateToken(user entities.User) (string, error) {
	response := network.Response{Body: &tokenSchema{}}
	request := network.Request{
		Path:   u.URL + "/tokens",
		Method: "POST",
		Body:   userSchema{Email: user.Email, Password: user.Password},
	}

	err := u.http.MakeRequest(request, &response, StatusErrors)
	if err != nil {
		return "", fmt.Errorf("error requesting for an user token: %w", err)
	}

	token := response.Body.(*tokenSchema)
	return token.Token, nil
}
