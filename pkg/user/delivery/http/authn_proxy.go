package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/user/entities"
)

// authn documentation: https://github.com/mainflux/mainflux/blob/master/authn/swagger.yaml

// AuthnProxy represents the interface to the authn's proxy operations
type AuthnProxy interface {
	CreateAppToken(user entities.User) (token string, err error)
}

// Authn is responsible for implementing the authn's proxy operations
type Authn struct {
	URL    string
	logger logging.Logger
}

// NewAuthnProxy creates a new authnProxy instance
func NewAuthnProxy(logger logging.Logger, authnHost string, authnPort uint16) *Authn {
	URL := fmt.Sprintf("http://%s:%d", authnHost, authnPort)
	logger.Debug("authn proxy configured to " + URL)
	return &Authn{URL, logger}
}

// authnRequest represents an authn service request
type authnRequest struct {
	Path          string
	Method        string
	Body          interface{}
	Authorization string
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
func (a *Authn) CreateAppToken(user entities.User) (string, error) {
	var response keySchema
	request := authnRequest{
		Path:          "/keys",
		Method:        "POST",
		Body:          keyRequestSchema{Issuer: user.Email, Type: 2, Duration: 31536000},
		Authorization: user.Token,
	}

	err := a.doRequest(request, &response)
	if err != nil {
		return "", fmt.Errorf("error requesting a new app token: %w", err)
	}

	return response.Value, nil
}

func (a *Authn) doRequest(request authnRequest, response interface{}) error {
	body, err := json.Marshal(&request.Body)
	if err != nil {
		return fmt.Errorf("error encoding body: %w", err)
	}

	req, err := http.NewRequest(request.Method, a.URL+request.Path, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error creating request object: %w", err)
	}
	req.Header.Add("Authorization", request.Authorization)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error executing HTTP request: %w", err)
	}
	defer resp.Body.Close()

	err = a.mapErrorFromStatusCode(resp.StatusCode)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(response)
	if err != nil {
		return fmt.Errorf("error decoding response: %w", err)
	}

	return nil
}

func (a *Authn) mapErrorFromStatusCode(code int) error {
	var err error

	switch code {
	case http.StatusBadRequest:
		err = entities.ErrMalformedRequest
	case http.StatusConflict:
		err = entities.ErrExistingID
	case http.StatusUnsupportedMediaType:
		err = entities.ErrMissingContentType
	case http.StatusInternalServerError:
		err = entities.ErrService
	}

	return err
}
