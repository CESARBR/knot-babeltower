package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
	"github.com/google/go-querystring/query"
)

// ThingProxy proxy a request to the thing service interface
type ThingProxy interface {
	Create(id, name, authorization string) (idGenerated string, err error)
}

type proxy struct {
	url    string
	logger logging.Logger
}

type objKnot struct {
	ID     string            `json:"id"`
	Schema []entities.Schema `json:"schema,omitempty"`
}

type objMetadata struct {
	Knot objKnot `json:"knot"`
}

// ThingProxyRepr is the entity that represents the thing on the remote thing's service
type ThingProxyRepr struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	Metadata objMetadata `json:"metadata"`
}

type pageFetchInput struct {
	Total  int               `json:"total"`
	Offset int               `json:"offset"`
	Limit  int               `json:"limit"`
	Things []*ThingProxyRepr `json:"things"`
}

// ErrThingNotFound represents the error when the schema has a invalid format
type ErrThingNotFound struct {
	ID string
}

func (etnf *ErrThingNotFound) Error() string {
	return fmt.Sprintf("Thing %s not found", etnf.ID)
}

func (p proxy) getJSONBody(id, name string, schemaList []entities.Schema) ([]byte, error) {
	body := ThingProxyRepr{
		Name: name,
		Metadata: objMetadata{
			Knot: objKnot{
				ID:     id,
				Schema: schemaList,
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

// RequestOptions represents the request query parameters
type RequestOptions struct {
	Limit  int `url:"limit"`
	Offset int `url:"offset"`
}

// RequestInfo aims to group all request releated information
type RequestInfo struct {
	method        string
	url           string
	authorization string
	contentType   string
	data          []byte
	options       *RequestOptions
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
	values, err := query.Values(info.options)
	if err != nil {
		return nil, err
	}
	queryString := "?" + values.Encode()

	/**
	 * Add Timeout in http.Client to avoid blocking the request.
	 */
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(info.method, info.url+queryString, bytes.NewBuffer(info.data))
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
	jsonBody, err := p.getJSONBody(id, name, nil)
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
		nil,
	}

	resp, err := p.sendRequest(requestInfo)
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
		nil,
	}

	resp, err := p.sendRequest(requestInfo)
	if err != nil {
		p.logger.Error(err)
		return err
	}

	defer resp.Body.Close()

	return p.mapErrorFromStatusCode(resp.StatusCode)
}

func (p proxy) getPaginatedThings(authorization string) ([]*ThingProxyRepr, error) {
	requestInfo := &RequestInfo{
		"GET",
		p.url + "/things",
		authorization,
		"application/json",
		nil,
		&RequestOptions{Limit: 100, Offset: 0}, // 100 is the max number of things that can be returned
	}

	var things []*ThingProxyRepr
	keepGoing := true
	for keepGoing {
		resp, err := p.sendRequest(requestInfo)
		if err != nil {
			p.logger.Error(err)
			return nil, err
		}
		defer resp.Body.Close()

		page := &pageFetchInput{}
		err = json.NewDecoder(resp.Body).Decode(&page)
		if err != nil {
			return nil, err
		}

		things = append(things, page.Things...)
		requestInfo.options.Offset += requestInfo.options.Limit

		if page.Total == len(things) {
			keepGoing = false
		}
	}

	return things, nil
}

func (p proxy) getThing(authorization, ID string) (*ThingProxyRepr, error) {
	things, err := p.getPaginatedThings(authorization)
	if err != nil {
		return nil, err
	}

	for i := range things {
		t := things[i]
		if t.Metadata.Knot.ID == ID {
			return t, nil
		}
	}

	return nil, &ErrThingNotFound{ID}
}
