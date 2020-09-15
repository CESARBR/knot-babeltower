package mocks

import (
	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
	"github.com/stretchr/testify/mock"
)

// FakeThingProxy represents a mocking type for the thing's proxy service
type FakeThingProxy struct {
	mock.Mock
	ReturnErr error
	CreateErr error
	Thing     *entities.Thing
}

// Create provides a mock function to create a thing on the thing's service
func (ftp *FakeThingProxy) Create(id, name, authorization string) (idGenerated string, err error) {
	ret := ftp.Called(id, name, authorization)
	return ret.String(0), ret.Error(1)
}

// UpdateConfig provides a mock function to update thing's config on the thing's service
func (ftp *FakeThingProxy) UpdateConfig(authorization, thingID string, config []entities.Config) error {
	ret := ftp.Called(authorization, thingID, config)
	return ret.Error(0)
}

// Get provides a mock function to receive a thing from the thing's service
func (ftp *FakeThingProxy) Get(authorization, thingID string) (*entities.Thing, error) {
	args := ftp.Called(authorization, thingID)
	return args.Get(0).(*entities.Thing), args.Error(1)
}

// List provides a mock function to list things from the thing's service
func (ftp *FakeThingProxy) List(authorization string) ([]*entities.Thing, error) {
	args := ftp.Called(authorization)
	return args.Get(0).([]*entities.Thing), args.Error(1)
}

// Remove provides a mock function to remove a thing on the thing's service
func (ftp *FakeThingProxy) Remove(authorization, thingID string) error {
	ret := ftp.Called(authorization, thingID)
	return ret.Error(0)
}
