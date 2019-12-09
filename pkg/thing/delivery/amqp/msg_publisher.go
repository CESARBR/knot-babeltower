package amqp

import (
	"encoding/json"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/network"
)

const (
	exchangeFogOut = "fogOut"
	registerOutKey = "device.registered"
)

// MsgPublisher handle messages received from a service
type MsgPublisher struct {
	logger logging.Logger
	amqp   *network.Amqp
}

// Publisher is the interface with methods that the publisher should have
type Publisher interface {
	SendRegisterDevice(network.RegisterResponseMsg) error
}

// NewMsgPublisher constructs the MsgPublisher
func NewMsgPublisher(logger logging.Logger, amqp *network.Amqp) *MsgPublisher {
	return &MsgPublisher{logger, amqp}
}

// SendRegisterDevice sends a registered message
func (mp *MsgPublisher) SendRegisterDevice(msg network.RegisterResponseMsg) error {
	mp.logger.Debug("Sending register message")

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		mp.logger.Error(err)
		return err
	}

	return mp.amqp.PublishPersistentMessage(exchangeFogOut, registerOutKey, jsonMsg)
}
