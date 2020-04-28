package mocks

import (
	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
	"github.com/stretchr/testify/mock"
)

// FakePublisher represents a mocking type for the publisher service
type FakePublisher struct {
	mock.Mock
	ReturnErr error
	SendError error
	Token     string
	Err       error
}

// SendRegisteredDevice provides a mock function to send a register device response
func (fp *FakePublisher) SendRegisteredDevice(thingID, token string, err error) error {
	ret := fp.Called(thingID, token, err)
	return ret.Error(0)
}

// SendUnregisteredDevice provides a mock function to send an unregister device response
func (fp *FakePublisher) SendUnregisteredDevice(thingID string, err error) error {
	ret := fp.Called(thingID, err)
	return ret.Error(0)
}

// SendUpdatedSchema provides a mock function to send an update schema response
func (fp *FakePublisher) SendUpdatedSchema(thingID string, err error) error {
	ret := fp.Called(thingID, err)
	return ret.Error(0)
}

// SendUpdateData provides a mock function to send an update data command
func (fp *FakePublisher) SendUpdateData(thingID string, data []entities.Data) error {
	args := fp.Called(thingID, data)
	return args.Error(0)
}

// SendRequestData provides a mock function to send a request data command
func (fp *FakePublisher) SendRequestData(thingID string, sensorIds []int) error {
	args := fp.Called(thingID, sensorIds)
	return args.Error(0)
}

// SendPublishedData provides a mock function to send a publish data command to connector
func (fp *FakePublisher) SendPublishedData(id string, data []entities.Data) error {
	ret := fp.Called(id, data)
	return ret.Error(0)
}
