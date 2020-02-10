package amqp

import (
	"encoding/json"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/network"
	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
)

const (
	exchangeConnIn    = "connIn"
	registerInKey     = "device.register"
	updateSchemaInKey = "schema.update"
)

type msgConnectorPublisher struct {
	logger logging.Logger
	amqp   *network.Amqp
}

// Connector handle messages received from a service
type ConnectorPublisher interface {
	SendRegisterDevice(string, string) error
	SendUpdateSchema(string, []entities.Schema) error
}

// NewConnector constructs the Connector
func NewMsgConnectorPublisher(logger logging.Logger, amqp *network.Amqp) ConnectorPublisher {
	return &msgConnectorPublisher{logger, amqp}
}

// SendRegisterDevice sends a registered message
func (mp *msgConnectorPublisher) SendRegisterDevice(id string, name string) error {
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
func (mp *msgConnectorPublisher) SendUpdateSchema(id string, schemaList []entities.Schema) error {
	mp.logger.Info("Sending update schema message to connector")
	msg := network.UpdateSchemaRequestMsg{ID: id, Schema: schemaList}
	bytes, err := json.Marshal(msg)
	if err != nil {
		mp.logger.Error(err)
		return err
	}
	return mp.amqp.PublishPersistentMessage(exchangeConnIn, updateSchemaInKey, bytes)
}
