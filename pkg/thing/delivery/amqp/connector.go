package amqp

import (
	"encoding/json"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/network"
	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
)

const (
	exchangeConnIn    = "connIn"
	exchangeConnOut   = "connOut"
	queueNameOut      = "connOut-messages"
	registerInKey     = "device.register"
	updateSchemaInKey = "schema.update"
)

type connector struct {
	logger logging.Logger
	amqp   *network.Amqp
}

// Connector handle messages received from a service
type Connector interface {
	SendRegisterDevice(string, string) error
	SendUpdateSchema(string, []entities.Schema) error
	RecvRegisterDevice() ([]byte, error)
}

// NewConnector constructs the Connector
func NewConnector(logger logging.Logger, amqp *network.Amqp) Connector {
	return &connector{logger, amqp}
}

// SendRegisterDevice sends a registered message
func (mp *connector) SendRegisterDevice(id string, name string) error {
	mp.logger.Debug("Sending register message")

	msg := network.RegisterRequestMsg{ID: id, Name: name}
	bytes, err := json.Marshal(msg)
	if err != nil {
		mp.logger.Error(err)
		return err
	}
	// TODO: receive message
	return mp.amqp.PublishPersistentMessage(exchangeConnIn, registerInKey, bytes)
}

// SendUpdateSchema sends an update schema message
func (mp *connector) SendUpdateSchema(id string, schemaList []entities.Schema) error {
	mp.logger.Info("Sending update schema message to connector")
	msg := network.UpdateSchemaRequestMsg{ID: id, Schema: schemaList}
	bytes, err := json.Marshal(msg)
	if err != nil {
		mp.logger.Error(err)
		return err
	}
	return mp.amqp.PublishPersistentMessage(exchangeConnIn, updateSchemaInKey, bytes)
}

// RecvRegisterDevice is a blocking function that receives the device
func (mp *connector) RecvRegisterDevice() ([]byte, error) {
	msgChan := make(chan network.InMsg)
	err := mp.amqp.OnMessage(msgChan, queueNameOut, exchangeConnOut, registerOutKey)
	if err != nil {
		mp.logger.Error(err)
		return nil, err
	}

	msg := <-msgChan
	mp.logger.Info("Message received:", string(msg.Body))

	return msg.Body, nil
}
