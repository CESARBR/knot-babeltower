package amqp

import (
	"encoding/json"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/network"
	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
)

const (
	exchangeFogOut    = "fogOut"
	registerOutKey    = "device.registered"
	schemaOutKey      = "schema.updated"
	listThingsOutKey  = "device.list"
	requestDataOutKey = "data.request"
)

// MsgPublisher handle messages received from a service
type MsgPublisher struct {
	logger logging.Logger
	amqp   *network.Amqp
}

// Publisher is the interface with methods that the publisher should have
type Publisher interface {
	SendRegisterDevice(network.RegisterResponseMsg) error
	SendUpdatedSchema(thingID string) error
	SendThings(things []*entities.Thing) error
	SendRequestData(thingID string, sensorIds []int) error
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

// SendUpdatedSchema sends the updated schema response
func (mp *MsgPublisher) SendUpdatedSchema(thingID string) error {
	resp := &network.UpdatedSchemaResponse{ID: thingID}
	msg, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	return mp.amqp.PublishPersistentMessage(exchangeFogOut, schemaOutKey, msg)
}

// SendThings sends the updated schema response
func (mp *MsgPublisher) SendThings(things []*entities.Thing) error {
	resp := &network.ListThingsResponse{Things: things}
	msg, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	return mp.amqp.PublishPersistentMessage(exchangeFogOut, listThingsOutKey, msg)
}

// SendRequestData sends request data command
func (mp *MsgPublisher) SendRequestData(thingID string, sensorIds []int) error {
	resp := &network.RequestDataCommand{ID: thingID, SensorIds: sensorIds}
	msg, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	return mp.amqp.PublishPersistentMessage(exchangeFogOut, requestDataOutKey, msg)
}
