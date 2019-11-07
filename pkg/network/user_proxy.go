package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/CESARBR/knot-babeltower/pkg/entities"
	"github.com/CESARBR/knot-babeltower/pkg/logging"
)

// UserProxy proxy a request to the user service interface
type UserProxy interface {
	Create(user entities.User) (err error)
}

// Proxy proxy a request to the user service
type Proxy struct {
	url    string
	logger logging.Logger
}

// NewUserProxy creates a proxy to the users service
func NewUserProxy(logger logging.Logger, hostname string, port uint16) *Proxy {
	url := fmt.Sprintf("http://%s:%d", hostname, port)

	logger.Debug("Proxy setup to " + url)
	return &Proxy{url, logger}
}

// Create proxy the http request to user service
func (p *Proxy) Create(user entities.User) (err error) {
	p.logger.Debug("Proxying request to create user")
	/**
	 * Add Timeout in http.Client to avoid blocking the request.
	 */
	client := &http.Client{Timeout: 10 * time.Second}
	jsonUser, err := json.Marshal(user)
	if err != nil {
		return err
	}

	resp, err := client.Post(p.url+"/users", "application/json", bytes.NewBuffer(jsonUser))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict {
		msg := fmt.Sprintf("User %s exists", user.Email)
		return entities.ErrEntityExists{Msg: msg}
	}

	return nil
}
