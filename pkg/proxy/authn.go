package proxy

import (
	"fmt"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/network"
	"github.com/CESARBR/knot-babeltower/pkg/user/entities"
)

// authn documentation: https://github.com/mainflux/mainflux/blob/master/authn/swagger.yaml

// AuthnProxy represents the interface to the authn's proxy operations
type AuthnProxy interface {
	CreateAppToken(user entities.User, duration int) (token string, err error)
}

// Authn is responsible for implementing the authn's proxy operations
type Authn struct {
	URL    string
	http   *network.HTTP
	logger logging.Logger
}

// NewAuthnProxy creates a new authnProxy instance
func NewAuthnProxy(logger logging.Logger, http *network.HTTP, authnHost string, authnPort uint16) AuthnProxy {
	URL := fmt.Sprintf("http://%s:%d", authnHost, authnPort)
	logger.Debug("authn proxy configured to " + URL)
	return &Authn{URL, http, logger}
}

// keyRequestSchema represents the schema for a key request
type keyRequestSchema struct {
	Issuer   string `json:"issuer"`
	Type     int    `json:"type"`
	Duration int    `json:"duration"`
}

// keySchema represents the schema for a key
type keySchema struct {
	ID        string `json:"id"`
	Value     string `json:"value"`
	IssuedAt  string `json:"issued_at"`
	ExpiresAt string `json:"expires_at"`
}

// CreateAppToken creates a valid token for the application
func (a *Authn) CreateAppToken(user entities.User, duration int) (string, error) {
	response := network.Response{Body: &keySchema{}}
	request := network.Request{
		Path:          a.URL + "/keys",
		Method:        "POST",
		Body:          keyRequestSchema{Issuer: user.Email, Type: 2, Duration: duration},
		Authorization: user.Token,
	}

	err := a.http.MakeRequest(request, &response, StatusErrors)
	if err != nil {
		return "", fmt.Errorf("error requesting a new app token: %w", err)
	}

	key := response.Body.(*keySchema)
	return key.Value, nil
}
