package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/user/entities"
)

// AuthnProxy is an interface to the Mainflux Authn service, which provides an API for.
// managing authentication keys. This interface provides a way to create application
// tokens (with configurable duration) and to validate the token used in operations across
// the platform.
// https://github.com/mainflux/mainflux/blob/0.11.0/authn/swagger.yaml
type AuthnProxy interface {
	CreateAppToken(user entities.User, duration int) (token string, err error)
}

// Authn takes a URL address that points to the Mainflux Authn service and implements
// the AuthnProxy interface methods.
type Authn struct {
	URL    string
	logger logging.Logger
}

// NewAuthnProxy creates a new authn instance and returns a pointer to the AuthnProxy interface
// implementation.
func NewAuthnProxy(logger logging.Logger, hostname string, port uint16) *Authn {
	URL := fmt.Sprintf("http://%s:%d", hostname, port)
	logger.Debug("authn proxy configured to " + URL)
	return &Authn{URL, logger}
}

// authnRequest represents a HTTP request to the Mainflux Authn service. It basically has
// fields that matches with the HTTP protocol structure (path, method, body, headers, etc).
type authnRequest struct {
	Path          string
	Method        string
	Body          interface{}
	Authorization string
}

// keyRequestSchema represents the request schema for creating a new token in the
// Mainflux platform.
// `Issuer` is the entity responsible for requesting the token creation.
// `Type` is the token type, which can be user (0) or app (2).
// `Duration` is the duration of the token until it expires (only for app tokens).
type keyRequestSchema struct {
	Issuer   string `json:"issuer"`
	Type     int    `json:"type"`
	Duration int    `json:"duration"`
}

// keyResponseSchema represents the response schema for the created token.
// `ID` is an unique ID that identifies the token.
// `Value` is the token value itself.
// `IssuedAt` and `ExpiresAt` are time fields to know when the token was created and
// when it expires.
type keyResponseSchema struct {
	ID        string `json:"id"`
	Value     string `json:"value"`
	IssuedAt  string `json:"issued_at"`
	ExpiresAt string `json:"expires_at"`
}

// CreateAppToken creates a new application token in the Mainflux platform. This type of
// token has a configurable duration.
func (a *Authn) CreateAppToken(user entities.User, duration int) (string, error) {
	var response keyResponseSchema
	request := authnRequest{
		Path:          "/keys",
		Method:        "POST",
		Body:          keyRequestSchema{Issuer: user.Email, Type: 2, Duration: duration},
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
