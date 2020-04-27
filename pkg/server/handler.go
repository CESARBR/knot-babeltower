package server

import (
	"errors"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/network"
	"github.com/CESARBR/knot-babeltower/pkg/thing/controllers"
)

// API definition to enable receiving request-reply commands from the clients
// The operations supported for this type of events are device authentication
// and list registered devices, as can be seen on the documentation:
// https://github.com/CESARBR/knot-babeltower/blob/master/docs/events.md
const (
	exchangeDevices       = "devices"
	exchangeDevicesType   = "direct"
	queueNameCommands     = "babeltower-command-messages"
	bindingKeyAuthDevice  = "device.auth"
	bindingKeyListDevices = "device.list"
)

const (
	queueFogIn               = "fogIn-messages"
	exchangeFogIn            = "fogIn"
	exchangeTypeFogIn        = "topic"
	queueConnOut             = "connOut-messages"
	exchangeConnOut          = "connOut"
	exchangeTypeConnOut      = "topic"
	bindingKeyDevice         = "device.*"
	bindingKeyPublishData    = "data.publish"
	bindingKeyDataCommands   = "data.*"
	bindingKeyDeviceCommands = "device.cmd.*"
	bindingKeySchema         = "schema.*"
)

// MsgHandler handle messages received from a service
type MsgHandler struct {
	logger          logging.Logger
	amqp            *network.Amqp
	thingController *controllers.ThingController
}

// NewMsgHandler creates a new MsgHandler instance with the necessary dependencies
func NewMsgHandler(logger logging.Logger, amqp *network.Amqp, thingController *controllers.ThingController) *MsgHandler {
	return &MsgHandler{logger, amqp, thingController}
}

// Start starts to listen messages
func (mc *MsgHandler) Start(started chan bool) {
	mc.logger.Debug("message handler started")
	msgChan := make(chan network.InMsg)
	err := mc.subscribeToMessages(msgChan)
	if err != nil {
		mc.logger.Error(err)
		started <- false
		return
	}

	go mc.onMsgReceived(msgChan)

	started <- true
}

// Stop stops to listen for messages
func (mc *MsgHandler) Stop() {
	mc.logger.Debug("message handler stopped")
}

func (mc *MsgHandler) subscribeToMessages(msgChan chan network.InMsg) error {
	var err error
	subscribe := func(msgChan chan network.InMsg, queue, exchange, kind, key string) {
		if err != nil {
			return
		}
		err = mc.amqp.OnMessage(msgChan, queue, exchange, kind, key)
	}

	// Subscribe to messages received from any client
	subscribe(msgChan, queueFogIn, exchangeFogIn, exchangeTypeFogIn, bindingKeyDevice)
	subscribe(msgChan, queueFogIn, exchangeFogIn, exchangeTypeFogIn, bindingKeySchema)
	subscribe(msgChan, queueFogIn, exchangeFogIn, exchangeTypeFogIn, bindingKeyDeviceCommands)
	subscribe(msgChan, queueFogIn, exchangeFogIn, exchangeTypeFogIn, bindingKeyPublishData)

	// Subscribe to messages received from the connector service
	subscribe(msgChan, queueConnOut, exchangeConnOut, exchangeTypeConnOut, bindingKeyDataCommands)
	subscribe(msgChan, queueConnOut, exchangeConnOut, exchangeTypeConnOut, bindingKeyDevice)

	// Subscribe to request-reply messages received from any client
	subscribe(msgChan, queueNameCommands, exchangeDevices, exchangeDevicesType, bindingKeyAuthDevice)
	subscribe(msgChan, queueNameCommands, exchangeDevices, exchangeDevicesType, bindingKeyListDevices)

	return err
}

func (mc *MsgHandler) onMsgReceived(msgChan chan network.InMsg) {
	for {
		var err error
		msg := <-msgChan
		mc.logger.Infof("exchange: %s, routing key: %s", msg.Exchange, msg.RoutingKey)
		mc.logger.Infof("message received: %s", string(msg.Body))

		token, ok := msg.Headers["Authorization"].(string)
		if !ok {
			mc.logger.Error(errors.New("authorization token not provided"))
			continue
		}

		if msg.Exchange == exchangeFogIn {
			err = mc.handleClientMessages(msg, token)
		} else if msg.Exchange == exchangeConnOut {
			err = mc.handleConnectorMessages(msg, token)
		} else if msg.Exchange == exchangeDevices {
			err = mc.handleCommandMessages(msg, token)
		}

		if err != nil {
			mc.logger.Error(err)
			continue
		}
	}
}

func (mc *MsgHandler) handleClientMessages(msg network.InMsg, token string) error {

	switch msg.RoutingKey {
	case "device.register":
		return mc.thingController.Register(msg.Body, token)
	case "device.unregister":
		return mc.thingController.Unregister(msg.Body, token)
	case "schema.update":
		return mc.thingController.UpdateSchema(msg.Body, token)
	case "data.publish":
		return mc.thingController.PublishData(msg.Body, token)
	}

	return nil
}

func (mc *MsgHandler) handleConnectorMessages(msg network.InMsg, token string) error {

	switch msg.RoutingKey {
	case "data.request":
		return mc.thingController.RequestData(msg.Body, token)
	case "data.update":
		return mc.thingController.UpdateData(msg.Body, token)
	case "device.registered":
		// Ignore message
	}

	return nil
}

func (mc *MsgHandler) handleCommandMessages(msg network.InMsg, token string) error {
	corrID, ok := msg.Headers["correlation-id"].(string)
	if !ok {
		return errors.New("correlation ID not provided")
	}

	switch msg.RoutingKey {
	case "device.auth":
		return mc.thingController.AuthDevice(msg.Body, token, corrID)
	case "device.list":
		return mc.thingController.ListDevices(token, corrID)
	}

	return nil
}
