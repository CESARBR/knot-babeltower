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
	exchangeDevices            = "device"
	exchangeDevicesType        = "direct"
	exchangeDataSentType       = "fanout"
	exchangeDataSent           = "data.sent"
	queueNameCommands          = "babeltower-command-messages"
	queueNameEvents            = "babeltower-event-messages"
	bindingKeyAuthDevice       = "device.auth"
	bindingKeyListDevices      = "device.list"
	bindingKeyRegisterDevice   = "device.register"
	bindingKeyUnregisterDevice = "device.unregister"
	bindingKeyRequestData      = "data.request"
	bindingKeyUpdateData       = "data.update"
	bindingKeySchemaSent       = "device.schema.sent"
	bindingKeyEmpty            = ""
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

	// Subscribe to general direct commands
	subscribe(msgChan, queueNameCommands, exchangeDevices, exchangeDevicesType, bindingKeyRegisterDevice)
	subscribe(msgChan, queueNameCommands, exchangeDevices, exchangeDevicesType, bindingKeyUnregisterDevice)
	subscribe(msgChan, queueNameCommands, exchangeDevices, exchangeDevicesType, bindingKeyRequestData)
	subscribe(msgChan, queueNameCommands, exchangeDevices, exchangeDevicesType, bindingKeyUpdateData)
	subscribe(msgChan, queueNameCommands, exchangeDevices, exchangeDevicesType, bindingKeySchemaSent)

	// Subscribe to request-reply messages received from any client
	subscribe(msgChan, queueNameCommands, exchangeDevices, exchangeDevicesType, bindingKeyAuthDevice)
	subscribe(msgChan, queueNameCommands, exchangeDevices, exchangeDevicesType, bindingKeyListDevices)

	// Subscribe to broadcasted data events
	subscribe(msgChan, queueNameEvents, exchangeDataSent, exchangeDataSentType, bindingKeyEmpty)

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

		if msg.RoutingKey == bindingKeyAuthDevice || msg.RoutingKey == bindingKeyListDevices {
			// handling request-reply command messages, which requires specific validations such as if reply_to was correctly received
			err = mc.handleRequestReplyCommands(msg, token)
		} else if msg.Exchange == exchangeDataSent {
			// handling broadcasted data events
			err = mc.handleBroadcastedData(msg, token)
		} else {
			// handling general direct commands
			err = mc.handleClientMessages(msg, token)
		}

		if err != nil {
			mc.logger.Error(err)
			continue
		}
	}
}

func (mc *MsgHandler) handleClientMessages(msg network.InMsg, token string) error {

	switch msg.RoutingKey {
	case bindingKeyRegisterDevice:
		return mc.thingController.Register(msg.Body, token)
	case bindingKeyUnregisterDevice:
		return mc.thingController.Unregister(msg.Body, token)
	case bindingKeySchemaSent:
		return mc.thingController.UpdateSchema(msg.Body, token)
	case bindingKeyRequestData:
		return mc.thingController.RequestData(msg.Body, token)
	case bindingKeyUpdateData:
		return mc.thingController.UpdateData(msg.Body, token)
	}

	return nil
}

func (mc *MsgHandler) handleRequestReplyCommands(msg network.InMsg, token string) error {
	replyTo, ok := msg.Headers["reply_to"].(string)
	if !ok {
		return errors.New("reply_to property not provided")
	}

	corrID, ok := msg.Headers["correlation_id"].(string)
	if !ok {
		return errors.New("correlation_id property not provided")
	}

	switch msg.RoutingKey {
	case bindingKeyAuthDevice:
		return mc.thingController.AuthDevice(msg.Body, token, replyTo, corrID)
	case bindingKeyListDevices:
		return mc.thingController.ListDevices(token, replyTo, corrID)
	}

	return nil
}

func (mc *MsgHandler) handleBroadcastedData(msg network.InMsg, token string) error {
	return mc.thingController.PublishData(msg.Body, token)
}
