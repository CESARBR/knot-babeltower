package amqp

import (
	"encoding/json"
	"fmt"

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
	updateDataOutKey  = "data.update"
	requestDataOutKey = "data.request"
)

// ClientPublisher is the interface with methods that the publisher should have
type ClientPublisher interface {
	SendRegisteredDevice(thingID, token string, err error) error
	SendUnregisteredDevice(thingID string, err error) error
	SendUpdatedSchema(thingID string, err error) error
	SendDevicesList(things []*entities.Thing, err error) error
	SendAuthStatus(thingID string, err error) error
	SendUpdateData(thingID string, data []entities.Data) error
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
func (mp *msgClientPublisher) SendRegisteredDevice(thingID, token string, err error) error {
	mp.logger.Debug("Sending registered message")
	errMsg := getErrMsg(err)
	resp := &network.DeviceRegisteredResponse{ID: thingID, Token: token, Error: errMsg}
	msg, err := json.Marshal(resp)
	if err != nil {
		mp.logger.Error(err)
		return err
	}

	return mp.amqp.PublishPersistentMessage(exchangeFogOut, registerOutKey, msg)
}

// SendUnregisterDevice publishes the unregistered device's id and error message to the device unregistered queue
func (mp *msgClientPublisher) SendUnregisteredDevice(thingID string, err error) error {
	mp.logger.Debug("Sending unregistered message")
	errMsg := getErrMsg(err)
	resp := &network.DeviceUnregisteredResponse{ID: thingID, Error: errMsg}
	msg, err := json.Marshal(resp)
	if err != nil {
		mp.logger.Error(err)
		return err
	}

	return mp.amqp.PublishPersistentMessage(exchangeFogOut, unregisterOutKey, msg)
}

// SendUpdatedSchema sends the updated schema response
func (mp *msgClientPublisher) SendUpdatedSchema(thingID string, err error) error {
	errMsg := getErrMsg(err)
	resp := &network.SchemaUpdatedResponse{ID: thingID, ErrMsg: errMsg}
	msg, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	return mp.amqp.PublishPersistentMessage(exchangeFogOut, schemaOutKey, msg)
}

// SendDevicesList sends the list devices command response
func (mp *msgClientPublisher) SendDevicesList(things []*entities.Thing, err error) error {
	errMsg := getErrMsg(err)
	resp := &network.DeviceListResponse{Things: things, ErrMsg: errMsg}
	msg, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	return mp.amqp.PublishPersistentMessage(exchangeFogOut, listThingsOutKey, msg)
}

// SendAuthStatus sends the auth thing status response
func (mp *msgClientPublisher) SendAuthStatus(thingID string, err error) error {
	errMsg := getErrMsg(err)
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

// SendUpdateData send update data command
func (mp *msgClientPublisher) SendUpdateData(thingID string, data []entities.Data) error {
	resp := &network.DataUpdate{ID: thingID, Data: data}
	msg, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("message parsing error: %w", err)
	}

	return mp.amqp.PublishPersistentMessage(exchangeFogOut, updateDataOutKey, msg)
}

func getErrMsg(err error) *string {
	if err != nil {
		msg := err.Error()
		return &msg
	}
	return nil
}
