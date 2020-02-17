package amqp

import (
	"encoding/json"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/network"
	"github.com/CESARBR/knot-babeltower/pkg/thing/interactors"
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
	thingInteractor interactors.Interactor
}

// NewMsgHandler constructs the MsgHandler
func NewMsgHandler(logger logging.Logger, amqp *network.Amqp, thingInteractor interactors.Interactor) *MsgHandler {
	return &MsgHandler{logger, amqp, thingInteractor}
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

func (mc *MsgHandler) handleRegisterMsg(body []byte, authorizationHeader string) error {
	msgParsed := network.RegisterRequestMsg{}
	err := json.Unmarshal(body, &msgParsed)
	if err != nil {
		return err
	}

	return mc.thingInteractor.Register(authorizationHeader, msgParsed.ID, msgParsed.Name)
}

func (mc *MsgHandler) handleUnregisterMsg(body []byte, authorizationHeader string) error {
	msg := network.UnregisterRequestMsg{}
	err := json.Unmarshal(body, &msg)
	if err != nil {
		return err
	}

	// TODO: call unregister device interactor
	return nil
}

func (mc *MsgHandler) handleUpdateSchemaMsg(body []byte, authorizationHeader string) error {
	var updateSchemaReq network.UpdateSchemaRequest
	err := json.Unmarshal(body, &updateSchemaReq)
	if err != nil {
		mc.logger.Error(err)
		return err
	}

	mc.logger.Info("Update schema message received")
	mc.logger.Debug(authorizationHeader, updateSchemaReq)

	err = mc.thingInteractor.UpdateSchema(authorizationHeader, updateSchemaReq.ID, updateSchemaReq.Schema)
	if err != nil {
		mc.logger.Error(err)
		return err
	}

	return nil
}

func (mc *MsgHandler) handleListDevices(authorization string) error {
	return mc.thingInteractor.List(authorization)
}

func (mc *MsgHandler) handleAuthDevice(body []byte, authorization string) error {
	var authThingReq network.AuthThingCommand
	err := json.Unmarshal(body, &authThingReq)
	if err != nil {
		mc.logger.Error(err)
		return err
	}

	mc.logger.Info("Auth device command received")
	mc.logger.Debug(authorization, authThingReq)
	return mc.thingInteractor.Auth(authorization, authThingReq.ID, authThingReq.Token)
}

func (mc *MsgHandler) handleRequestData(body []byte, authorization string) error {
	var requestDataReq network.RequestDataCommand
	err := json.Unmarshal(body, &requestDataReq)
	if err != nil {
		mc.logger.Error(err)
		return err
	}

	mc.logger.Info("Request data command received")
	mc.logger.Debug(authorization, requestDataReq)
	err = mc.thingInteractor.RequestData(authorization, requestDataReq.ID, requestDataReq.SensorIds)
	if err != nil {
		return err
	}

	return nil
}

func (mc *MsgHandler) onMsgReceived(msgChan chan network.InMsg) {
	for {
		msg := <-msgChan
		mc.logger.Infof("Exchange: %s, routing key: %s", msg.Exchange, msg.RoutingKey)
		mc.logger.Infof("Message received: %s", string(msg.Body))

		authorizationHeader := msg.Headers["Authorization"]

		switch msg.RoutingKey {
		case "device.register":
			err := mc.handleRegisterMsg(msg.Body, authorizationHeader.(string))
			if err != nil {
				mc.logger.Error(err)
				continue
			}
		case "device.registered":
			// Ignore message
			continue
		case "device.unregister":
			err := mc.handleUnregisterMsg(msg.Body, authorizationHeader.(string))
			if err != nil {
				mc.logger.Error(err)
				continue
			}
		case "device.cmd.list":
			mc.logger.Info("List things request received")
			err := mc.handleListDevices(authorizationHeader.(string))
			if err != nil {
				mc.logger.Error(err)
				continue
			}
		case "device.cmd.auth":
			err := mc.handleAuthDevice(msg.Body, authorizationHeader.(string))
			if err != nil {
				mc.logger.Error(err)
				continue
			}
		case "data.request":
			err := mc.handleRequestData(msg.Body, authorizationHeader.(string))
			if err != nil {
				mc.logger.Error(err)
				continue
			}
		case "schema.update":
			err := mc.handleUpdateSchemaMsg(msg.Body, authorizationHeader.(string))
			if err != nil {
				mc.logger.Error(err)
				continue
			}
		}
	}
}
