package mocks

import (
	"github.com/CESARBR/knot-babeltower/pkg/network"
	"github.com/stretchr/testify/mock"
)

// FakeAmqpReceiver represents a mock type for amqp receiver
type FakeAmqpReceiver struct {
	mock.Mock
}

// OnMessage receives the message
func (f *FakeAmqpReceiver) OnMessage(msgChan chan network.InMsg, queueName string, exchangeName string, exchangeType string, key string) error {
	args := f.Called(msgChan, queueName, exchangeName, exchangeType, key)
	return args.Error(0)
}
