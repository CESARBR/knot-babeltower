package network

import (
	"encoding/json"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
)

const (
	exchangeConnIn = "connIn"
	registerInKey  = "device.register"
)

type connector struct {
	logger logging.Logger
	amqp   *Amqp
}

// Connector handle messages received from a service
type Connector interface {
	SendRegisterDevice(string, string) error
}

// NewConnector constructs the Connector
func NewConnector(logger logging.Logger, amqp *Amqp) Connector {
	return &connector{logger, amqp}
}

// SendRegisterDevice sends a registered message
func (mp *connector) SendRegisterDevice(id string, name string) error {
	mp.logger.Debug("Sending register message")

	msg := RegisterRequestMsg{id, name}
	bytes, err := json.Marshal(msg)
	if err != nil {
		mp.logger.Error(err)
		return err
	}
	// TODO: receive message
	return mp.amqp.PublishPersistentMessage(exchangeConnIn, registerInKey, bytes)
}
