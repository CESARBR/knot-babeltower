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
	unregisterOutKey  = "device.unregistered"
	schemaOutKey      = "schema.updated"
	listThingsOutKey  = "device.list"
	authDeviceOutKey  = "device.auth"
	requestDataOutKey = "data.request"
)

// ClientPublisher is the interface with methods that the publisher should have
type ClientPublisher interface {
	SendRegisteredDevice(network.DeviceRegisteredResponse) error
	SendUnregisteredDevice(thingID string, errMsg *string) error
	SendUpdatedSchema(thingID string) error
	SendThings(things []*entities.Thing) error
	SendAuthStatus(thingID string, errMsg *string) error
	SendRequestData(thingID string, sensorIds []int) error
}

// msgClientPublisher handle messages received from a service
type msgClientPublisher struct {
	logger logging.Logger
	amqp   *network.Amqp
}

// NewMsgClientPublisher constructs the msgClientPublisher
func NewMsgClientPublisher(logger logging.Logger, amqp *network.Amqp) ClientPublisher {
	return &msgClientPublisher{logger, amqp}
}

// SendRegisterDevice publishes the registered device's credentials to the device registration queue
func (mp *msgClientPublisher) SendRegisteredDevice(msg network.DeviceRegisteredResponse) error {
	mp.logger.Debug("Sending registered message")
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		mp.logger.Error(err)
		return err
	}

	return mp.amqp.PublishPersistentMessage(exchangeFogOut, registerOutKey, jsonMsg)
}

// SendUnregisterDevice publishes the unregistered device's id and error message to the device unregistered queue
func (mp *msgClientPublisher) SendUnregisteredDevice(thingID string, errMsg *string) error {
	mp.logger.Debug("Sending unregistered message")
	resp := &network.DeviceUnregisteredResponse{ID: thingID, Error: errMsg}
	msg, err := json.Marshal(resp)
	if err != nil {
		mp.logger.Error(err)
		return err
	}

	return mp.amqp.PublishPersistentMessage(exchangeFogOut, unregisterOutKey, msg)
}

// SendUpdatedSchema sends the updated schema response
func (mp *msgClientPublisher) SendUpdatedSchema(thingID string) error {
	resp := &network.SchemaUpdatedResponse{ID: thingID}
	msg, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	return mp.amqp.PublishPersistentMessage(exchangeFogOut, schemaOutKey, msg)
}

// SendThings sends the updated schema response
func (mp *msgClientPublisher) SendThings(things []*entities.Thing) error {
	resp := &network.DeviceListResponse{Things: things}
	msg, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	return mp.amqp.PublishPersistentMessage(exchangeFogOut, listThingsOutKey, msg)
}

// SendAuthStatus sends the auth thing status response
func (mp *msgClientPublisher) SendAuthStatus(thingID string, errMsg *string) error {
	resp := &network.DeviceAuthResponse{ID: thingID, ErrMsg: errMsg}
	msg, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	return mp.amqp.PublishPersistentMessage(exchangeFogOut, authDeviceOutKey, msg)
}

// SendRequestData sends request data command
func (mp *msgClientPublisher) SendRequestData(thingID string, sensorIds []int) error {
	resp := &network.DataRequest{ID: thingID, SensorIds: sensorIds}
	msg, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	return mp.amqp.PublishPersistentMessage(exchangeFogOut, requestDataOutKey, msg)
}
