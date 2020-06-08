package amqp

import (
	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/network"
	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
)

const (
	exchangeDevice            = "device"
	exchangeDeviceType        = "direct"
	exchangeDataPublished     = "data.published"
	exchangeDataPublishedType = "fanout"
	registerOutKey            = "device.registered"
	unregisterOutKey          = "device.unregistered"
	schemaOutKey              = "device.schema.updated"
	updateDataKey             = "data.update"
	requestDataKey            = "data.request"
)

// Publisher provides methods to send events to the clients
type Publisher interface {
	PublishRegisteredDevice(thingID, name, token string, err error) error
	PublishUnregisteredDevice(thingID string, err error) error
	PublishUpdatedSchema(thingID string, schema []entities.Schema, err error) error
	PublishUpdateData(thingID string, data []entities.Data) error
	PublishRequestData(thingID string, sensorIds []int) error
	PublishPublishedData(thingID, token string, data []entities.Data) error
}

// Sender represents the operations to send commands response
type Sender interface {
	SendAuthResponse(thingID, replyTo, corrID string, err error) error
	SendListResponse(things []*entities.Thing, replyTo, corrID string, err error) error
}

// msgClientPublisher handle messages received from a service
type msgClientPublisher struct {
	logger logging.Logger
	amqp   *network.Amqp
}

// commandSender handle messages received from a service
type commandSender struct {
	logger logging.Logger
	amqp   *network.Amqp
}

// NewMsgClientPublisher constructs the msgClientPublisher
func NewMsgClientPublisher(logger logging.Logger, amqp *network.Amqp) Publisher {
	return &msgClientPublisher{logger, amqp}
}

// NewCommandSender creates a new commandSender instance
func NewCommandSender(logger logging.Logger, amqp *network.Amqp) Sender {
	return &commandSender{logger, amqp}
}

// PublishRegisteredDevice publishes the registered device's credentials to the device registration queue
func (mp *msgClientPublisher) PublishRegisteredDevice(thingID, name, token string, err error) error {
	mp.logger.Debug("sending registered response")
	errMsg := getErrMsg(err)
	msg := network.NewMessage(network.DeviceRegisteredResponse{ID: thingID, Name: name, Token: token, Error: errMsg})

	return mp.amqp.PublishPersistentMessage(exchangeDevice, exchangeDeviceType, registerOutKey, msg, nil)
}

// PublishUnregisteredDevice publishes the unregistered device's id and error message to the device unregistered queue
func (mp *msgClientPublisher) PublishUnregisteredDevice(thingID string, err error) error {
	mp.logger.Debug("sending unregistered response")
	errMsg := getErrMsg(err)
	msg := network.NewMessage(network.DeviceUnregisteredResponse{ID: thingID, Error: errMsg})

	return mp.amqp.PublishPersistentMessage(exchangeDevice, exchangeDeviceType, unregisterOutKey, msg, nil)
}

// PublishUpdatedSchema sends the updated schema response
func (mp *msgClientPublisher) PublishUpdatedSchema(thingID string, schema []entities.Schema, err error) error {
	mp.logger.Debug("sending update schema response")
	errMsg := getErrMsg(err)
	msg := network.NewMessage(network.SchemaUpdatedResponse{ID: thingID, Schema: schema, Error: errMsg})

	return mp.amqp.PublishPersistentMessage(exchangeDevice, exchangeDeviceType, schemaOutKey, msg, nil)
}

// PublishRequestData sends request data command
func (mp *msgClientPublisher) PublishRequestData(thingID string, sensorIds []int) error {
	mp.logger.Debug("sending request data request")
	msg := network.NewMessage(network.DataRequest{ID: thingID, SensorIds: sensorIds})
	routingKey := "device." + thingID + "." + requestDataKey

	return mp.amqp.PublishPersistentMessage(exchangeDevice, exchangeDeviceType, routingKey, msg, nil)
}

// PublishUpdateData send update data command
func (mp *msgClientPublisher) PublishUpdateData(thingID string, data []entities.Data) error {
	mp.logger.Debug("sending update data request")
	msg := network.NewMessage(network.DataUpdate{ID: thingID, Data: data})
	routingKey := "device." + thingID + "." + updateDataKey

	return mp.amqp.PublishPersistentMessage(exchangeDevice, exchangeDeviceType, routingKey, msg, nil)
}

// SendAuthResponse sends the auth thing status response
func (cs *commandSender) SendAuthResponse(thingID string, replyTo, corrID string, err error) error {
	cs.logger.Debug("sending auth device response")
	errMsg := getErrMsg(err)
	msg := network.NewMessage(network.DeviceAuthResponse{ID: thingID, Error: errMsg})
	options := &network.MessageOptions{CorrelationID: corrID}

	return cs.amqp.PublishPersistentMessage(exchangeDevice, exchangeDeviceType, replyTo, msg, options)
}

// SendListResponse sends the list devices command response
func (cs *commandSender) SendListResponse(things []*entities.Thing, replyTo, corrID string, err error) error {
	cs.logger.Debug("sending list devices response")
	errMsg := getErrMsg(err)
	msg := network.NewMessage(network.DeviceListResponse{Things: things, Error: errMsg})
	options := &network.MessageOptions{CorrelationID: corrID}

	return cs.amqp.PublishPersistentMessage(exchangeDevice, exchangeDeviceType, replyTo, msg, options)
}

// PublishPublishedData send update data command
func (mp *msgClientPublisher) PublishPublishedData(thingID, token string, data []entities.Data) error {
	mp.logger.Debug("sending publish data request")
	msg := network.NewMessage(network.DataSent{ID: thingID, Data: data})
	options := &network.MessageOptions{Authorization: token}

	return mp.amqp.PublishPersistentMessage(exchangeDataPublished, exchangeDataPublishedType, "", msg, options)
}

func getErrMsg(err error) *string {
	if err != nil {
		msg := err.Error()
		return &msg
	}
	return nil
}
