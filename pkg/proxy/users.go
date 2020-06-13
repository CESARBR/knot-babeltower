package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/user/entities"
)

// UsersProxy represents the interface to the user's proxy operations
type UsersProxy interface {
	Create(user entities.User) (err error)
	CreateToken(user entities.User) (token string, err error)
}

// Users is responsible for implementing the user's proxy operations
type Users struct {
	URL    string
	logger logging.Logger
}

// UserTokenResponse represents the create token response from the user's service
type UserTokenResponse struct {
	Token string `json:"token"`
}

// NewUsersProxy creates a new Proxy instance
func NewUsersProxy(logger logging.Logger, userHost string, userPort uint16) *Users {
	URL := fmt.Sprintf("http://%s:%d", userHost, userPort)
	logger.Debug("user proxy configured to " + URL)
	return &Users{URL, logger}
}

// Create proxy the http request to user service
func (p *Users) Create(user entities.User) (err error) {
	p.logger.Debug("proxying request to create user")
	/**
	 * Add Timeout in http.Client to avoid blocking the request.
	 */
	client := &http.Client{Timeout: 10 * time.Second}
	jsonUser, err := json.Marshal(user)
	if err != nil {
		return err
	}

	resp, err := client.Post(p.URL+"/users", "application/json", bytes.NewBuffer(jsonUser))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return p.mapErrorFromStatusCode(resp.StatusCode)
}

// CreateToken creates a valid token for the specified user
func (p *Users) CreateToken(user entities.User) (string, error) {
	var resp *http.Response

	credentials, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err = client.Post(p.URL+"/tokens", "application/json", bytes.NewBuffer(credentials))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	err = p.mapErrorFromStatusCode(resp.StatusCode)
	if err != nil {
		return "", err
	}

	tr := &UserTokenResponse{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(tr)
	if err != nil {
		return "", nil
	}

	return tr.Token, nil
}

func (p *Users) mapErrorFromStatusCode(code int) error {
	var err error

	if code != http.StatusCreated {
		switch code {
		case http.StatusForbidden:
			err = entities.ErrUserForbidden
		case http.StatusConflict:
			err = entities.ErrUserExists
		case http.StatusBadRequest:
			err = entities.ErrUserBadRequest
		}
	}

	return err
}
