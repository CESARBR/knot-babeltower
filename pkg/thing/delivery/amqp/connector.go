package amqp

import (
	"encoding/json"
	"fmt"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/network"
	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
)

const (
	exchangeConnIn    = "connIn"
	registerInKey     = "device.register"
	unregisterInKey   = "device.unregister"
	updateSchemaInKey = "schema.update"
	publishDataInKey  = "data.publish"
)

// ConnectorPublisher handle messages received from a service
type ConnectorPublisher interface {
	SendRegisterDevice(string, string) error
	SendUnregisterDevice(string) error
	SendUpdateSchema(string, []entities.Schema) error
	SendPublishData(string, []entities.Data) error
}

type msgConnectorPublisher struct {
	logger logging.Logger
	amqp   *network.Amqp
}

// NewMsgConnectorPublisher constructs the Connector
func NewMsgConnectorPublisher(logger logging.Logger, amqp *network.Amqp) ConnectorPublisher {
	return &msgConnectorPublisher{logger, amqp}
}

// SendRegisterDevice sends a register message
func (mp *msgConnectorPublisher) SendRegisterDevice(id string, name string) error {
	mp.logger.Debug("sending register message")

	msg := network.DeviceRegisterRequest{ID: id, Name: name}
	bytes, err := json.Marshal(msg)
	if err != nil {
		mp.logger.Error(err)
		return err
	}
	// TODO: receive message
	return mp.amqp.PublishPersistentMessage(exchangeConnIn, registerInKey, bytes)
}

// SendUnregisterDevice sends an unregister message
func (mp *msgConnectorPublisher) SendUnregisterDevice(id string) error {
	mp.logger.Debug("sending unregister message")
	msg := network.DeviceUnregisterRequest{ID: id}
	bytes, err := json.Marshal(msg)
	if err != nil {
		mp.logger.Error(err)
		return err
	}

	return mp.amqp.PublishPersistentMessage(exchangeConnIn, unregisterInKey, bytes)
}

// SendUpdateSchema sends an update schema message
func (mp *msgConnectorPublisher) SendUpdateSchema(id string, schemaList []entities.Schema) error {
	mp.logger.Info("sending update schema message to connector")
	msg := network.SchemaUpdateRequest{ID: id, Schema: schemaList}
	bytes, err := json.Marshal(msg)
	if err != nil {
		mp.logger.Error(err)
		return err
	}
	return mp.amqp.PublishPersistentMessage(exchangeConnIn, updateSchemaInKey, bytes)
}

// SendPublishData sends a publish data message
func (mp *msgConnectorPublisher) SendPublishData(id string, data []entities.Data) error {
	mp.logger.Info("sending publish data message to connector")
	msg := network.DataPublish{ID: id, Data: data}
	bytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("message parsing error: %w", err)
	}
	return mp.amqp.PublishPersistentMessage(exchangeConnIn, publishDataInKey, bytes)
}
