package server

import (
	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/network"
	"github.com/CESARBR/knot-babeltower/pkg/thing/controllers"
)

const (
	queueNameFogIn           = "fogIn-messages"
	exchangeFogIn            = "fogIn"
	queueNameConnOut         = "connOut-messages"
	exchangeConnOut          = "connOut"
	bindingKeyDevice         = "device.*"
	bindingKeyData           = "data.*"
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
	subscribe := func(msgChan chan network.InMsg, queueName, exchange, key string) {
		if err != nil {
			return
		}
		err = mc.amqp.OnMessage(msgChan, queueName, exchange, key)
	}

	// Subscribe to messages received from any client
	subscribe(msgChan, queueNameFogIn, exchangeFogIn, bindingKeyDevice)
	subscribe(msgChan, queueNameFogIn, exchangeFogIn, bindingKeySchema)
	subscribe(msgChan, queueNameFogIn, exchangeFogIn, bindingKeyDeviceCommands)

	// Subscribe to messages received from the connector service
	subscribe(msgChan, queueNameConnOut, exchangeConnOut, bindingKeyData)
	subscribe(msgChan, queueNameConnOut, exchangeConnOut, bindingKeyDevice)

	return err
}

func (mc *MsgHandler) onMsgReceived(msgChan chan network.InMsg) {
	for {
		var err error
		msg := <-msgChan
		mc.logger.Infof("exchange: %s, routing key: %s", msg.Exchange, msg.RoutingKey)
		mc.logger.Infof("message received: %s", string(msg.Body))

		if msg.Exchange == exchangeFogIn {
			err = mc.handleClientMessages(msg)
		} else if msg.Exchange == exchangeConnOut {
			err = mc.handleConnectorMessages(msg)
		}

		if err != nil {
			mc.logger.Error(err)
			continue
		}
	}
}

func (mc *MsgHandler) handleClientMessages(msg network.InMsg) error {
	authorizationHeader := msg.Headers["Authorization"]

	switch msg.RoutingKey {
	case "device.register":
		return mc.thingController.Register(msg.Body, authorizationHeader.(string))
	case "device.unregister":
		return mc.thingController.Unregister(msg.Body, authorizationHeader.(string))
	case "schema.update":
		return mc.thingController.UpdateSchema(msg.Body, authorizationHeader.(string))
	case "device.cmd.auth":
		return mc.thingController.AuthDevice(msg.Body, authorizationHeader.(string))
	case "device.cmd.list":
		return mc.thingController.ListDevices(authorizationHeader.(string))
	}

	return nil
}

func (mc *MsgHandler) handleConnectorMessages(msg network.InMsg) error {
	authorizationHeader := msg.Headers["Authorization"]

	switch msg.RoutingKey {
	case "data.request":
		return mc.thingController.RequestData(msg.Body, authorizationHeader.(string))
	case "data.update":
		return mc.thingController.UpdateData(msg.Body, authorizationHeader.(string))
	case "data.publish":
		return mc.thingController.PublishData(msg.Body, authorizationHeader.(string))
	case "device.registered":
		// Ignore message
	}

	return nil
}
