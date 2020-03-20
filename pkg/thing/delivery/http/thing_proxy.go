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

type errorConflict struct{ error }

func (err errorConflict) Error() string {
	return "Error conflict"
}

// ThingProxy proxy a request to the thing service interface
type ThingProxy interface {
	Create(id, name, authorization string) (idGenerated string, err error)
	UpdateSchema(authorization, ID string, schemaList []entities.Schema) error
	List(authorization string) (things []*entities.Thing, err error)
	Get(authorization, ID string) (*entities.Thing, error)
	Remove(authorization, ID string) error
}

// ThingProxyRepr is the entity that represents the thing on the remote thing's service
type ThingProxyRepr struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	Metadata objMetadata `json:"metadata"`
}

type objMetadata struct {
	Knot objKnot `json:"knot"`
}

type objKnot struct {
	ID     string            `json:"id"`
	Schema []entities.Schema `json:"schema,omitempty"`
}

type pageFetchInput struct {
	Total  int               `json:"total"`
	Offset int               `json:"offset"`
	Limit  int               `json:"limit"`
	Things []*ThingProxyRepr `json:"things"`
}

type proxy struct {
	url    string
	logger logging.Logger
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

// RequestOptions represents the request query parameters
type RequestOptions struct {
	Limit  int `url:"limit"`
	Offset int `url:"offset"`
}

// NewThingProxy creates a proxy to the thing service
func NewThingProxy(logger logging.Logger, hostname string, port uint16) ThingProxy {
	url := fmt.Sprintf("http://%s:%d", hostname, port)

	logger.Debug("Proxy setup to " + url)
	return proxy{url, logger}
}

// Create register a thing on service and return the id generated
func (p proxy) Create(id, name, authorization string) (idGenerated string, err error) {
	p.logger.Debug("Proxying request to create thing")
	t := p.getRemoteThingRepr(id, name, nil)
	body, err := json.Marshal(t)
	if err != nil {
		p.logger.Error(err)
		return "", err
	}

	requestInfo := &RequestInfo{
		"POST",
		p.url + "/things",
		authorization,
		"application/json",
		body,
		nil,
	}

	resp, err := p.sendRequest(requestInfo)
	if err != nil {
		p.logger.Error(err)
		return "", err
	}
	defer resp.Body.Close()

	err = p.mapErrorFromStatusCode(resp.StatusCode)
	if err != nil {
		p.logger.Error(err)
		return "", err
	}

	locationHeader := resp.Header.Get("Location")
	thingID := locationHeader[len("/things/"):] // get substring after "/things/"
	return thingID, nil
}

// UpdateSchema receives the thing's ID and schema and send a HTTP request to
// the thing's service in order to update it with the schema.
func (p proxy) UpdateSchema(authorization, ID string, schemaList []entities.Schema) error {
	t, err := p.Get(authorization, ID)
	if err != nil {
		return err
	}

	rt := p.getRemoteThingRepr(t.ID, t.Name, t.Schema)
	rt.Metadata.Knot.Schema = schemaList
	parsedBody, err := json.Marshal(rt)
	if err != nil {
		p.logger.Error(err)
		return err
	}

	requestInfo := &RequestInfo{
		"PUT",
		p.url + "/things/" + t.Token,
		authorization,
		"application/json",
		parsedBody,
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

func (p proxy) List(authorization string) (things []*entities.Thing, err error) {
	pagThings, err := p.getPaginatedThings(authorization)
	if err != nil {
		return nil, nil
	}

	for _, t := range pagThings {
		things = append(things, &entities.Thing{ID: t.Metadata.Knot.ID, Name: t.Name, Schema: t.Metadata.Knot.Schema})
	}

	return things, err
}

// Get list the things registered on thing's service
func (p proxy) Get(authorization, ID string) (*entities.Thing, error) {
	things, err := p.getPaginatedThings(authorization)
	if err != nil {
		return nil, err
	}

	for i := range things {
		t := things[i]
		if t.Metadata.Knot.ID == ID {
			nt := &entities.Thing{ID: ID, Token: t.ID, Name: t.Name, Schema: t.Metadata.Knot.Schema}
			return nt, nil
		}
	}

	return nil, entities.ErrThingNotFound
}

// Remove removes the indicated thing from the thing's service
func (p proxy) Remove(authorization, ID string) error {
	t, err := p.Get(authorization, ID)
	if err != nil {
		return err
	}

	requestInfo := &RequestInfo{
		"DELETE",
		p.url + "/things/" + t.Token,
		authorization,
		"application/json",
		nil,
		nil,
	}

	resp, err := p.sendRequest(requestInfo)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return p.mapErrorFromStatusCode(resp.StatusCode)
}

func (p proxy) getRemoteThingRepr(id, name string, schemaList []entities.Schema) ThingProxyRepr {
	return ThingProxyRepr{
		Name: name,
		Metadata: objMetadata{
			Knot: objKnot{
				ID:     id,
				Schema: schemaList,
			},
		},
	}
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

func (p proxy) mapErrorFromStatusCode(code int) error {
	var err error

	if code != http.StatusCreated {
		switch code {
		case http.StatusConflict:
			err = errorConflict{}
		case http.StatusForbidden:
			err = entities.ErrThingForbidden
		}
	}
	return err
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

		err = p.mapErrorFromStatusCode(resp.StatusCode)
		if err != nil {
			p.logger.Error(err)
			return nil, err
		}

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
