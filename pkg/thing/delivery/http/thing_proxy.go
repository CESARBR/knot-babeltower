package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
)

// ThingProxy proxy a request to the thing service interface
type ThingProxy interface {
	Create(id, name, authorization string) (idGenerated string, err error)
	UpdateSchema(ID string, schemaList []entities.Schema) error
}

type proxy struct {
	url    string
	logger logging.Logger
}

type objKnot struct {
	ID string `json:"id"`
}

type objMetadata struct {
	Knot objKnot `json:"knot"`
}

func (p proxy) getJSONBody(id, name string) ([]byte, error) {
	body := struct {
		Name     string      `json:"name"`
		Metadata objMetadata `json:"metadata"`
	}{
		Name: name,
		Metadata: objMetadata{
			Knot: objKnot{
				ID: id,
			},
		},
	}
	return json.Marshal(body)
}

// NewThingProxy creates a proxy to the thing service
func NewThingProxy(logger logging.Logger, hostname string, port uint16) ThingProxy {
	url := fmt.Sprintf("http://%s:%d", hostname, port)

	logger.Debug("Proxy setup to " + url)
	return proxy{url, logger}
}

// RequestInfo aims to group all request releated information
type RequestInfo struct {
	method        string
	url           string
	authorization string
	contentType   string
	data          []byte
}

type errorConflict struct{ error }
type errorForbidden struct{ error }

func (err errorForbidden) Error() string {
	return "Error forbidden"
}

func (err errorConflict) Error() string {
	return "Error conflict"
}

func (p proxy) mapErrorFromStatusCode(code int) error {
	var err error

	if code != http.StatusCreated {
		switch code {
		case http.StatusConflict:
			err = errorConflict{}
		case http.StatusForbidden:
			err = errorForbidden{}
		}
	}
	return err
}

func (p proxy) sendRequest(info *RequestInfo) (*http.Response, error) {
	/**
	 * Add Timeout in http.Client to avoid blocking the request.
	 */
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(info.method, info.url, bytes.NewBuffer(info.data))
	if err != nil {
		p.logger.Error(err)
		return nil, err
	}

	req.Header.Set("Authorization", info.authorization)
	req.Header.Set("Content-Type", info.contentType)

	return client.Do(req)
}

// Create register a thing on service and return the id generated
func (p proxy) Create(id, name, authorization string) (idGenerated string, err error) {
	p.logger.Debug("Proxying request to create thing")
	jsonBody, err := p.getJSONBody(id, name)
	if err != nil {
		p.logger.Error(err)
		return "", err
	}

	requestInfo := &RequestInfo{
		"POST",
		p.url + "/things",
		authorization,
		"application/json",
		jsonBody,
	}

	resp, err := p.sendRequest(requestInfo)
	print(resp)
	if err != nil {
		p.logger.Error(err)
		return "", err
	}

	locationHeader := resp.Header.Get("Location")
	fmt.Print(locationHeader)
	thingID := locationHeader[len("/things/"):] // get substring after "/things/"
	return thingID, p.mapErrorFromStatusCode(resp.StatusCode)
}

// UpdateSchema receives the thing's ID and schema and send a HTTP request to
// the thing's service in order to update it with the schema.
func (p proxy) UpdateSchema(ID string, schemaList []entities.Schema) error {
	parsedSchema, err := json.Marshal(schemaList)
	if err != nil {
		p.logger.Error(err)
		return err
	}

	requestInfo := &RequestInfo{
		"PUT",
		p.url + "/things" + ID,
		"authorization",
		"application/json",
		parsedSchema,
	}

	resp, err := p.sendRequest(requestInfo)
	if err != nil {
		p.logger.Error(err)
		return err
	}

	defer resp.Body.Close()

	return p.mapErrorFromStatusCode(resp.StatusCode)
}
