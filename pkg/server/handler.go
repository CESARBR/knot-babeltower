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

// NewMsgHandler constructs the MsgHandler
func NewMsgHandler(logger logging.Logger, amqp *network.Amqp, thingController *controllers.ThingController) *MsgHandler {
	return &MsgHandler{logger, amqp, thingController}
}

// Start starts to listen messages
func (mc *MsgHandler) Start(started chan bool) {
	mc.logger.Debug("Msg handler started")
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
	mc.logger.Debug("Msg handler stopped")
}

func (mc *MsgHandler) subscribeToMessages(msgChan chan network.InMsg) error {
	var err error
	subscribe := func(msgChan chan network.InMsg, queueName, exchange, key string) {
		if err != nil {
			return
		}
		err = mc.amqp.OnMessage(msgChan, queueName, exchange, key)
	}

	subscribe(msgChan, queueNameFogIn, exchangeFogIn, bindingKeyDevice)
	subscribe(msgChan, queueNameFogIn, exchangeFogIn, bindingKeySchema)
	subscribe(msgChan, queueNameFogIn, exchangeFogIn, bindingKeyDeviceCommands)
	subscribe(msgChan, queueNameConnOut, exchangeConnOut, bindingKeyData)
	subscribe(msgChan, queueNameConnOut, exchangeConnOut, bindingKeyDevice)
	return err
}

func (mc *MsgHandler) onMsgReceived(msgChan chan network.InMsg) {
	for {
		msg := <-msgChan
		mc.logger.Infof("Exchange: %s, routing key: %s", msg.Exchange, msg.RoutingKey)
		mc.logger.Infof("Message received: %s", string(msg.Body))

		authorizationHeader := msg.Headers["Authorization"]

		switch msg.RoutingKey {
		case "device.register":
			err := mc.thingController.Register(msg.Body, authorizationHeader.(string))
			if err != nil {
				mc.logger.Error(err)
				continue
			}
		case "device.registered":
			// Ignore message
			continue
		case "device.unregister":
			err := mc.thingController.Unregister(msg.Body, authorizationHeader.(string))
			if err != nil {
				mc.logger.Error(err)
				continue
			}
		case "device.cmd.list":
			mc.logger.Info("List things request received")
			err := mc.thingController.ListDevices(authorizationHeader.(string))
			if err != nil {
				mc.logger.Error(err)
				continue
			}
		case "device.cmd.auth":
			err := mc.thingController.AuthDevice(msg.Body, authorizationHeader.(string))
			if err != nil {
				mc.logger.Error(err)
				continue
			}
		case "data.request":
			err := mc.thingController.RequestData(msg.Body, authorizationHeader.(string))
			if err != nil {
				mc.logger.Error(err)
				continue
			}
		case "schema.update":
			err := mc.thingController.UpdateSchema(msg.Body, authorizationHeader.(string))
			if err != nil {
				mc.logger.Error(err)
				continue
			}
		}
	}
}
