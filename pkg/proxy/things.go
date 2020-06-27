package proxy

import (
	"fmt"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/network"
	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
)

// things documentation: https://github.com/mainflux/mainflux/blob/master/things/swagger.yaml

// ThingProxy proxy a request to the thing service interface
type ThingProxy interface {
	Create(id, name, authorization string) (idGenerated string, err error)
	UpdateSchema(authorization, ID string, schemaList []entities.Schema) error
	UpdateConfig(authorization, ID string, configList []entities.Config) error
	List(authorization string) (things []*entities.Thing, err error)
	Get(authorization, ID string) (*entities.Thing, error)
	Remove(authorization, ID string) error
}

type proxy struct {
	url    string
	http   *network.HTTP
	logger logging.Logger
}

// NewThingProxy creates a proxy to the thing service
func NewThingProxy(logger logging.Logger, http *network.HTTP, hostname string, port uint16) ThingProxy {
	url := fmt.Sprintf("http://%s:%d", hostname, port)
	logger.Debug("proxy setup to " + url)
	return proxy{url, http, logger}
}

// createThingReqSchema represents the schema for a create thing request
type createThingReqSchema struct {
	Key      string      `json:"key,omitempty"`
	Name     string      `json:"name"`
	Metadata interface{} `json:"metadata"`
}

// updateThingReqSchema represents the schema for an update thing request
type updateThingReqSchema struct {
	Name     string      `json:"name,omitempty"`
	Metadata interface{} `json:"metadata,omitempty"`
}

// thingsPageSchema represents the schema for a page of things
type thingsPageSchema struct {
	Things []*thingResSchema `json:"things"`
	Total  int               `json:"total"`
	Offset int               `json:"offset"`
	Limit  int               `json:"limit"`
}

// thingResSchema represents the schema for a thing
type thingResSchema struct {
	ID       string       `json:"id"`
	Key      string       `json:"key"`
	Name     string       `json:"name"`
	Metadata metadataKNoT `json:"metadata"`
}

// thingsQuery represents the query parameters for requests to things service
type thingsQuery struct {
	Limit    int    `url:"limit,omitempty"`
	Offset   int    `url:"offset,omitempty"`
	Name     string `url:"name,omitempty"`
	Metadata string `url:"metadata,omitempty"`
}

// metadataKNoT represents the KNoT metadata
type metadataKNoT struct {
	KNoT thingKNoT `json:"knot"`
}

// thingKNoT represents a thing KNoT on things service
type thingKNoT struct {
	ID     string            `json:"id"`
	Schema []entities.Schema `json:"schema,omitempty"`
	Config []entities.Config `json:"config,omitempty"`
}

// Create register a thing on service and return the id generated
func (p proxy) Create(id, name, authorization string) (string, error) {
	var response network.Response
	request := network.Request{
		Path:          p.url + "/things",
		Method:        "POST",
		Body:          createThingReqSchema{Name: name, Metadata: metadataKNoT{KNoT: thingKNoT{ID: id}}},
		Authorization: authorization,
	}

	err := p.http.MakeRequest(request, &response, StatusErrors)
	if err != nil {
		return "", fmt.Errorf("error creating a new thing: %w", err)
	}

	location := response.Header.Get("Location")
	return location[len("/things/"):], nil
}

// UpdateSchema receives the thing's ID and schema and send a HTTP request to
// the thing's service in order to update it with the schema.
func (p proxy) UpdateSchema(authorization, ID string, schemaList []entities.Schema) error {
	thing, err := p.Get(authorization, ID)
	if err != nil {
		return fmt.Errorf("error getting thing: %w", err)
	}

	request := network.Request{
		Path:          p.url + "/things/" + thing.Token,
		Method:        "PUT",
		Body:          updateThingReqSchema{Metadata: metadataKNoT{KNoT: thingKNoT{ID: ID, Schema: schemaList, Config: thing.Config}}},
		Authorization: authorization,
	}

	err = p.http.MakeRequest(request, nil, StatusErrors)
	if err != nil {
		return fmt.Errorf("error requesting for update thing: %w", err)
	}

	return nil
}

// UpdateConfig receives as parameters the authorization token, thing's ID and config. After that,
// it sends a HTTP request to the thing's service in order to update it with the new config.
func (p proxy) UpdateConfig(authorization, ID string, configList []entities.Config) error {
	thing, err := p.Get(authorization, ID)
	if err != nil {
		return fmt.Errorf("error getting thing: %w", err)
	}

	request := network.Request{
		Path:          p.url + "/things/" + thing.Token,
		Method:        "PUT",
		Body:          updateThingReqSchema{Metadata: metadataKNoT{KNoT: thingKNoT{ID: ID, Schema: thing.Schema, Config: configList}}},
		Authorization: authorization,
	}

	err = p.http.MakeRequest(request, nil, StatusErrors)
	if err != nil {
		return fmt.Errorf("error requesting for update thing: %w", err)
	}

	return nil
}

func (p proxy) List(authorization string) ([]*entities.Thing, error) {
	things, err := p.getPaginatedThings(authorization)
	if err != nil {
		return nil, fmt.Errorf("error listing things: %w", err)
	}

	return things, nil
}

// Get list the things registered on thing's service
func (p proxy) Get(authorization, ID string) (*entities.Thing, error) {
	things, err := p.getPaginatedThings(authorization)
	if err != nil {
		return nil, fmt.Errorf("error getting things: %w", err)
	}

	for _, t := range things {
		if t.ID == ID {
			return t, nil
		}
	}

	return nil, entities.ErrThingNotFound
}

// Remove removes the indicated thing from the thing's service
func (p proxy) Remove(authorization, ID string) error {
	thing, err := p.Get(authorization, ID)
	if err != nil {
		return fmt.Errorf("error getting thing: %w", err)
	}

	request := network.Request{
		Path:          p.url + "/things/" + thing.Token,
		Method:        "DELETE",
		Authorization: authorization,
	}

	err = p.http.MakeRequest(request, nil, StatusErrors)
	if err != nil {
		return fmt.Errorf("error requesting to delete thing: %w", err)
	}

	return nil
}

func (p proxy) getPaginatedThings(authorization string) ([]*entities.Thing, error) {
	paginatedThings := []*entities.Thing{}

	for offset := 0; offset == len(paginatedThings); offset += 100 {
		things, err := p.getThings(authorization, offset)
		if err != nil {
			return nil, fmt.Errorf("error getting paginated things: %w", err)
		}

		paginatedThings = append(paginatedThings, things...)
	}

	return paginatedThings, nil
}

func (p proxy) getThings(authorization string, offset int) ([]*entities.Thing, error) {
	response := network.Response{Body: &thingsPageSchema{}}
	request := network.Request{
		Path:          p.url + "/things",
		Method:        "GET",
		Query:         thingsQuery{Limit: 100, Offset: offset},
		Authorization: authorization,
	}

	err := p.http.MakeRequest(request, &response, StatusErrors)
	if err != nil {
		return nil, fmt.Errorf("error requesting for things: %w", err)
	}

	things := []*entities.Thing{}
	for _, t := range response.Body.(*thingsPageSchema).Things {
		things = append(things, &entities.Thing{
			ID:     t.Metadata.KNoT.ID,
			Token:  t.ID,
			Name:   t.Name,
			Schema: t.Metadata.KNoT.Schema,
			Config: t.Metadata.KNoT.Config,
		})
	}

	return things, nil
}
