package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/user/entities"
)

// HTTP handles HTTP requests
type HTTP struct {
	logger logging.Logger
}

// Request represents a thing service request
type Request struct {
	Path          string
	Method        string
	Body          interface{}
	Authorization string
}

// NewHTTP constructs the HTTP requests handler
func NewHTTP(logger logging.Logger) *HTTP {
	return &HTTP{logger}
}

// MakeRequest execute a HTTP request
func (h *HTTP) MakeRequest(request Request, response interface{}) error {
	body, err := json.Marshal(&request.Body)
	if err != nil {
		return fmt.Errorf("error encoding body: %w", err)
	}

	req, err := http.NewRequest(request.Method, request.Path, bytes.NewBuffer(body))
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

	err = h.mapErrorFromStatusCode(resp.StatusCode)
	if err != nil {
		return err
	}

	if response != nil {
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(response)
		if err != nil {
			return fmt.Errorf("error decoding response: %w", err)
		}
	}

	return nil
}

func (h *HTTP) mapErrorFromStatusCode(code int) error {
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
