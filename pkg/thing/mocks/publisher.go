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
