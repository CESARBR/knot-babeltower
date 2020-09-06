package proxy

import (
	"fmt"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/network"
	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
)

// things documentation: https://github.com/mainflux/mainflux/blob/master/things/swagger.yaml

// ThingsProxy proxy a request to the thing service interface
type ThingsProxy interface {
	Create(ID, name, authorization string) (token string, err error)
	UpdateSchema(authorization, ID string, schemaList []entities.Schema) (err error)
	UpdateConfig(authorizaiton, ID string, configList []entities.Config) (err error)
	List(authorization string) (things []*entities.Thing, err error)
	Get(authorization, ID string) (thing *entities.Thing, err error)
	Remove(authorization, ID string) (err error)
}

type things struct {
	URL    string
	http   *network.HTTP
	logger logging.Logger
}

// NewThingsProxy creates a proxy to the thing service
func NewThingsProxy(logger logging.Logger, http *network.HTTP, hostname string, port uint16) ThingsProxy {
	URL := fmt.Sprintf("http://%s:%d", hostname, port)
	logger.Debug("things proxy configured to " + URL)
	return &things{URL, http, logger}
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

const pageSize = 100

// Create register a thing on service and return the id generated
func (t *things) Create(ID, name, authorization string) (string, error) {
	var response network.Response
	request := network.Request{
		Path:          t.URL + "/things",
		Method:        "POST",
		Body:          createThingReqSchema{Name: name, Metadata: metadataKNoT{KNoT: thingKNoT{ID: ID}}},
		Authorization: authorization,
	}

	err := t.http.MakeRequest(request, &response, StatusErrors)
	if err != nil {
		return "", fmt.Errorf("error creating a new thing: %w", err)
	}

	location := response.Header.Get("Location")
	return location[len("/things/"):], nil
}

// UpdateSchema receives the thing's ID and schema and send a HTTP request to
// the thing's service in order to update it with the schema.
func (t *things) UpdateSchema(authorization, ID string, schemaList []entities.Schema) error {
	thing, err := t.Get(authorization, ID)
	if err != nil {
		return fmt.Errorf("error getting thing: %w", err)
	}

	request := network.Request{
		Path:          t.URL + "/things/" + thing.Token,
		Method:        "PUT",
		Body:          updateThingReqSchema{Metadata: metadataKNoT{KNoT: thingKNoT{ID: ID, Schema: schemaList, Config: thing.Config}}},
		Authorization: authorization,
	}

	err = t.http.MakeRequest(request, nil, StatusErrors)
	if err != nil {
		return fmt.Errorf("error requesting for update thing: %w", err)
	}

	return nil
}

// UpdateConfig receives as parameters the authorization token, thing's ID and config. After that,
// it sends a HTTP request to the thing's service in order to update it with the new config.
func (t *things) UpdateConfig(authorization, ID string, configList []entities.Config) error {
	thing, err := t.Get(authorization, ID)
	if err != nil {
		return fmt.Errorf("error getting thing: %w", err)
	}

	request := network.Request{
		Path:          t.URL + "/things/" + thing.Token,
		Method:        "PUT",
		Body:          updateThingReqSchema{Metadata: metadataKNoT{KNoT: thingKNoT{ID: ID, Schema: thing.Schema, Config: configList}}},
		Authorization: authorization,
	}

	err = t.http.MakeRequest(request, nil, StatusErrors)
	if err != nil {
		return fmt.Errorf("error requesting for update thing: %w", err)
	}

	return nil
}

func (t *things) List(authorization string) ([]*entities.Thing, error) {
	things, err := t.getPaginatedThings(authorization)
	if err != nil {
		return nil, fmt.Errorf("error listing things: %w", err)
	}

	return things, nil
}

// Get list the things registered on thing's service
func (t *things) Get(authorization, ID string) (*entities.Thing, error) {
	things, err := t.getPaginatedThings(authorization)
	if err != nil {
		return nil, fmt.Errorf("error getting things: %w", err)
	}

	for _, th := range things {
		if th.ID == ID {
			return th, nil
		}
	}

	return nil, entities.ErrThingNotFound
}

// Remove removes the indicated thing from the thing's service
func (t *things) Remove(authorization, ID string) error {
	thing, err := t.Get(authorization, ID)
	if err != nil {
		return fmt.Errorf("error getting thing: %w", err)
	}

	request := network.Request{
		Path:          t.URL + "/things/" + thing.Token,
		Method:        "DELETE",
		Authorization: authorization,
	}

	err = t.http.MakeRequest(request, nil, StatusErrors)
	if err != nil {
		return fmt.Errorf("error requesting to delete thing: %w", err)
	}

	return nil
}

func (t *things) getPaginatedThings(authorization string) ([]*entities.Thing, error) {
	paginatedThings := []*entities.Thing{}

	for offset := 0; offset == len(paginatedThings); offset += pageSize {
		things, err := t.getThings(authorization, offset)
		if err != nil {
			return nil, fmt.Errorf("error getting paginated things: %w", err)
		}

		paginatedThings = append(paginatedThings, things...)
	}

	return paginatedThings, nil
}

func (t *things) getThings(authorization string, offset int) ([]*entities.Thing, error) {
	response := network.Response{Body: &thingsPageSchema{}}
	request := network.Request{
		Path:          t.URL + "/things",
		Method:        "GET",
		Query:         thingsQuery{Limit: pageSize, Offset: offset},
		Authorization: authorization,
	}

	err := t.http.MakeRequest(request, &response, StatusErrors)
	if err != nil {
		return nil, fmt.Errorf("error requesting for things: %w", err)
	}

	things := []*entities.Thing{}
	for _, th := range response.Body.(*thingsPageSchema).Things {
		things = append(things, &entities.Thing{
			ID:     th.Metadata.KNoT.ID,
			Token:  th.ID,
			Name:   th.Name,
			Schema: th.Metadata.KNoT.Schema,
			Config: th.Metadata.KNoT.Config,
		})
	}

	return things, nil
}
