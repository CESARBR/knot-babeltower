package mocks

import (
	"github.com/CESARBR/knot-babeltower/pkg/network"
	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
	"github.com/stretchr/testify/mock"
)

// FakePublisher represents a mocking type for the publisher service
type FakePublisher struct {
	mock.Mock
	ReturnErr error
	SendError error
	Token     string
	ErrMsg    *string
}

// SendRegisterDevice provides a mock function to send a register device response
func (fp *FakePublisher) SendRegisterDevice(msg network.RegisterResponseMsg) error {
	ret := fp.Called(msg)
	return ret.Error(0)
}

// SendUpdatedSchema provides a mock function to send an update schema response
func (fp *FakePublisher) SendUpdatedSchema(thingID string) error {
	ret := fp.Called(thingID)
	return ret.Error(0)
}

// SendThings provides a mock function to send a list things response
func (fp *FakePublisher) SendThings(things []*entities.Thing) error {
	args := fp.Called(things)
	return args.Error(0)
}

// SendAuthStatus provides a mock function to send auth thing command response
func (fp *FakePublisher) SendAuthStatus(thingID string, errMsg *string) error {
	args := fp.Called(thingID, errMsg)
	return args.Error(0)
}

// SendRequestData provides a mock function to send a request data command
func (fp *FakePublisher) SendRequestData(thingID string, sensorIds []int) error {
	args := fp.Called(thingID, sensorIds)
	return args.Error(0)
}
