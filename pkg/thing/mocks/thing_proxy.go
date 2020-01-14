package mocks

import (
	"github.com/CESARBR/knot-babeltower/pkg/thing/delivery/http"
	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
	"github.com/stretchr/testify/mock"
)

// FakeThingProxy represents a mocking type for the thing's proxy service
type FakeThingProxy struct {
	mock.Mock
	ReturnErr error
}

// Create provides a mock function to create a thing on the thing's seervice
func (ftp *FakeThingProxy) Create(id, name, authorization string) (idGenerated string, err error) {
	ret := ftp.Called(id, name, authorization)
	return ret.String(0), ret.Error(1)
}

// UpdateSchema provides a mock function to update thing's schema on the thing's service
func (ftp *FakeThingProxy) UpdateSchema(authorization, thingID string, schema []entities.Schema) error {
	ret := ftp.Called(thingID, schema)
	return ret.Error(0)
}

// Get provides a mock function to receive a thing from the thing's service
func (ftp *FakeThingProxy) Get(authorization, thingID string) (*http.ThingProxyRepr, error) {
	return nil, nil
}

// List provides a mock function to list things from the thing's service
func (ftp *FakeThingProxy) List(authorization string) ([]*entities.Thing, error) {
	return nil, nil
}
