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

// PublishRegisteredDevice provides a mock function to send a register device response
func (fp *FakePublisher) PublishRegisteredDevice(thingID, name, token string, err error) error {
	ret := fp.Called(thingID, name, token, err)
	return ret.Error(0)
}

// PublishUnregisteredDevice provides a mock function to send an unregister device response
func (fp *FakePublisher) PublishUnregisteredDevice(thingID, token string, err error) error {
	ret := fp.Called(thingID, err)
	return ret.Error(0)
}

// PublishUpdatedConfig provides a mock function to send an update config response
func (fp *FakePublisher) PublishUpdatedConfig(thingID string, config []entities.Config, changed bool, err error) error {
	ret := fp.Called(thingID, config, changed, err)
	return ret.Error(0)
}

// PublishUpdateData provides a mock function to send an update data command
func (fp *FakePublisher) PublishUpdateData(thingID string, data []entities.Data) error {
	args := fp.Called(thingID, data)
	return args.Error(0)
}

// PublishRequestData provides a mock function to send a request data command
func (fp *FakePublisher) PublishRequestData(thingID string, sensorIds []int) error {
	args := fp.Called(thingID, sensorIds)
	return args.Error(0)
}

// PublishBroadcastData provides a mock function to publish data in broadcast mode
func (fp *FakePublisher) PublishBroadcastData(thingID, token string, data []entities.Data) error {
	args := fp.Called(thingID, token, data)
	return args.Error(0)
}
