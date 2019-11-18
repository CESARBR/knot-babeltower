package network

import (
	"encoding/json"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
)

const (
	exchangeFogOut = "FogOut"
	registerOutKey = "device.registered"
)

// MsgPublisher handle messages received from a service
type MsgPublisher struct {
	logger logging.Logger
	amqp   *Amqp
}

// NewMsgPublisher constructs the MsgPublisher
func NewMsgPublisher(logger logging.Logger, amqp *Amqp) *MsgPublisher {
	return &MsgPublisher{logger, amqp}
}

// SendRegisterDevice sends a registered message
func (mp *MsgPublisher) SendRegisterDevice(msg RegisterResponseMsg) error {
	mp.logger.Debug("Sending register message")

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		mp.logger.Error(err)
		return err
	}

	return mp.amqp.PublishPersistentMessage(exchangeFogOut, registerOutKey, jsonMsg)
}
