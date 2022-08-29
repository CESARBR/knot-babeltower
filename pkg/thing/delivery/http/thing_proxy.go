package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
	"github.com/google/go-querystring/query"
)

const maxNumberThingsThatCanBeReturned = 100

type errorConflict struct{ error }

func (err errorConflict) Error() string {
	return "error conflict"
}

// ThingProxy is an interface to the Mainflux Things service, which provides an API for
// managing things (logical representation of a physical thing in IoT). This interface
// provides a set of operations to manage things (CRUD). In addition, it supports the
// updating of thing's configuration by using the Mainflux metadata capabilities.
// https://github.com/mainflux/mainflux/blob/0.12.1/things/openapi.yml
type ThingProxy interface {
	Create(id, name, authorization string) (string, error)
	UpdateConfig(authorization, ID string, configList []entities.Config) error
	List(authorization string) (things []*entities.Thing, err error)
	Get(authorization, ID string) (*entities.Thing, error)
	Remove(authorization, ID string) error
}

type thingProxy struct {
	url       string
	logger    logging.Logger
	cache     map[string]*entities.Thing
	cacheLock *sync.Mutex
}

type pageFetchSchema struct {
	Total  int            `json:"total"`
	Offset int            `json:"offset"`
	Limit  int            `json:"limit"`
	Things []*thingSchema `json:"things"`
}

type thingSchema struct {
	ID       string        `json:"id"`
	Name     string        `json:"name"`
	Metadata thingMetadata `json:"metadata"`
}

type thingMetadata struct {
	Thing knotThing `json:"knot"`
}

type knotThing struct {
	ID     string            `json:"id"`
	Config []entities.Config `json:"config,omitempty"`
}

type requestInfo struct {
	method        string
	url           string
	authorization string
	contentType   string
	data          []byte
	options       *requestOptions
}

type requestOptions struct {
	Limit  int `url:"limit"`
	Offset int `url:"offset"`
}

// NewThingProxy creates a new things instance and returns a pointer to the ThingsProxy interface
// implementation.
func NewThingProxy(logger logging.Logger, hostname, protocol string, port uint16) *thingProxy {
	url := fmt.Sprintf("%s://%s:%d", protocol, hostname, port)
	logger.Debug("things proxy configured to " + url)
	cache := make(map[string]*entities.Thing)
	cacheLock := &sync.Mutex{}
	return &thingProxy{url, logger, cache, cacheLock}
}

// Create registers a new thing in the Mainflux platform. It receives the thing's properties and
// map them to the Mainflux internal representation. As a result, the operation returns the things
// ID.
func (p thingProxy) Create(id, name, authorization string) (string, error) {
	p.logger.Info("Creating thing")
	t := p.getThingSchema(id, name, nil)
	body, err := json.Marshal(t)
	if err != nil {
		return "", err
	}

	requestInfo := &requestInfo{
		"POST",
		p.url + "/things",
		authorization,
		"application/json",
		body,
		nil,
	}

	resp, err := p.sendRequest(requestInfo)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	err = p.mapErrorFromStatusCode(resp.StatusCode)
	if err != nil {
		return "", err
	}

	location := resp.Header.Get("Location")
	return location[len("/things/"):], nil // get substring after "/things/"
}

// UpdateConfig updates the internal thing's representation with the config in the format supported
// by the KNoT protocol. KNoT Thing config has two data structures: (1) schema and (2) event.
// (1) represents the sensor semantic models (temperature, voltage, etc).
// (2) represents the sensor data publishing configuration (interval, custom behavior when data changes, etc).
func (p thingProxy) UpdateConfig(authorization, ID string, configList []entities.Config) error {
	p.logger.Info("Updating config")
	t, err := p.Get(authorization, ID)
	if err != nil {
		return err
	}

	rt := p.getThingSchema(t.ID, t.Name, t.Config)
	rt.Metadata.Thing.Config = configList
	parsedBody, err := json.Marshal(rt)
	if err != nil {
		return err
	}

	requestInfo := &requestInfo{
		"PUT",
		p.url + "/things/" + t.Token,
		authorization,
		"application/json",
		parsedBody,
		nil,
	}

	resp, err := p.sendRequest(requestInfo)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	updatedThing := thingFactory(t.ID, t.Token, t.Name, configList)
	updateCache(p.cache, ID, updatedThing, p.cacheLock)
	return p.mapErrorFromStatusCode(resp.StatusCode)
}

func updateCache(cache map[string]*entities.Thing, ID string, thing *entities.Thing, lock *sync.Mutex) {
	lock.Lock()
	cache[ID] = thing
	lock.Unlock()
}

// List returns the registered things according to the KNoT Cloud representation.
// The Mainflux Things API blocks requests for a large number of things. Thus,
// this method paginates over them a returns a single slice of things.
func (p thingProxy) List(authorization string) ([]*entities.Thing, error) {
	p.logger.Info("Listing things")
	things := []*entities.Thing{}
	pagThings, err := p.getPaginatedThings(authorization)
	if err != nil {
		return things, err
	}

	for _, t := range pagThings {
		things = append(things, &entities.Thing{ID: t.Metadata.Thing.ID, Name: t.Name, Config: t.Metadata.Thing.Config})
	}

	return things, err
}

// Get retrieves an invidual thing from the Mainflux service. It uses the KNoT Thing's ID as filter.
func (p thingProxy) Get(authorization, ID string) (*entities.Thing, error) {
	cachedThing, ok := getThingFromCache(p.cache, ID, p.cacheLock)
	if ok {
		return cachedThing, nil
	} else {
		p.logger.Infof("Not cached thing ID: %s\n", ID)
		things, err := p.getPaginatedThings(authorization)
		if err != nil {
			return nil, err
		}

		specifiedThing := findSpecifiedThing(things, ID)
		if specifiedThing != nil {
			nt := thingFactory(ID, specifiedThing.ID, specifiedThing.Name, specifiedThing.Metadata.Thing.Config)
			updateCache(p.cache, ID, nt, p.cacheLock)
			return nt, nil
		}
		return nil, entities.ErrThingNotFound
	}
}
func findSpecifiedThing(things []*thingSchema, ID string) *thingSchema {
	for _, thing := range things {
		if thing.Metadata.Thing.ID == ID {
			return thing
		}
	}
	return nil
}

func thingFactory(ID string, token string, name string, config []entities.Config) *entities.Thing {
	return &entities.Thing{ID: ID, Token: token, Name: name, Config: config}
}

func getThingFromCache(cache map[string]*entities.Thing, ID string, lock *sync.Mutex) (*entities.Thing, bool) {
	lock.Lock()
	thing, ok := cache[ID]
	lock.Unlock()
	return thing, ok
}

// Remove removes an individual thing from the Mainflux service. It uses the KNoT Thing's ID as filter.
func (p thingProxy) Remove(authorization, ID string) error {
	p.logger.Info("Removing thing")
	t, err := p.Get(authorization, ID)
	if err != nil {
		return err
	}

	requestInfo := &requestInfo{
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
	removeThingFromCache(p.cache, ID, p.cacheLock)
	return p.mapErrorFromStatusCode(resp.StatusCode)
}

func removeThingFromCache(cache map[string]*entities.Thing, ID string, lock *sync.Mutex) {
	lock.Lock()
	delete(cache, ID)
	lock.Unlock()
}

func (p thingProxy) getThingSchema(id, name string, configList []entities.Config) thingSchema {
	return thingSchema{
		Name: name,
		Metadata: thingMetadata{
			Thing: knotThing{
				ID:     id,
				Config: configList,
			},
		},
	}
}

func (p thingProxy) getPaginatedThings(authorization string) ([]*thingSchema, error) {
	requestInfo := &requestInfo{
		"GET",
		p.url + "/things",
		authorization,
		"application/json",
		nil,
		&requestOptions{Limit: maxNumberThingsThatCanBeReturned, Offset: 0},
	}

	var things []*thingSchema
	keepGoing := true
	for keepGoing {
		resp, err := p.sendRequest(requestInfo)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		err = p.mapErrorFromStatusCode(resp.StatusCode)
		if err != nil {
			return nil, err
		}

		page := &pageFetchSchema{}
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

func (p thingProxy) sendRequest(info *requestInfo) (*http.Response, error) {
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
		return nil, err
	}

	req.Header.Set("Authorization", info.authorization)
	req.Header.Set("Content-Type", info.contentType)

	return client.Do(req)
}

func (p thingProxy) mapErrorFromStatusCode(code int) error {
	var err error
	errorMapping := httpCodeErrorMappingFactory()
	if code != http.StatusCreated {
		err = errorMapping[code]
	}
	return err
}

func httpCodeErrorMappingFactory() map[int]error {
	const conflictHTTPCode = 409
	const unauthorizedHTTPCode = 401
	errorMapping := make(map[int]error)
	errorMapping[conflictHTTPCode] = errorConflict{}
	errorMapping[unauthorizedHTTPCode] = entities.ErrThingUnauthorized
	return errorMapping
}
