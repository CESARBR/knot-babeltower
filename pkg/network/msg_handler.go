package network

import (
	"encoding/json"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
)

const (
	queueNameFogIn  = "fogIn-messages"
	exchangeFogIn   = "fogIn"
	bindingKeyFogIn = "device.*"
)

// MsgHandler handle messages received from a service
type MsgHandler struct {
	logger        logging.Logger
	amqp          *Amqp
	registerThing Interactor
}

// Interactor is the use case to be executed
type Interactor interface {
	Execute(id string, args ...interface{}) error
}

func (mc *MsgHandler) handleRegisterMsg(body []byte, authorizationHeader string) error {
	msgParsed := RegisterRequestMsg{}
	err := json.Unmarshal(body, &msgParsed)
	if err != nil {
		return err
	}

	return mc.registerThing.Execute(msgParsed.ID, msgParsed.Name, authorizationHeader)
}

func (mc *MsgHandler) onMsgReceived(msgChan chan InMsg) {
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
		}
	}
}

// NewMsgHandler constructs the MsgHandler
func NewMsgHandler(logger logging.Logger, amqp *Amqp, registerThing Interactor) *MsgHandler {
	return &MsgHandler{logger, amqp, registerThing}
}

// Start starts to listen for messages
func (mc *MsgHandler) Start(started chan bool) {
	mc.logger.Debug("Msg handler started")

	msgChan := make(chan InMsg)
	err := mc.amqp.OnMessage(msgChan, queueNameFogIn, exchangeFogIn, bindingKeyFogIn)
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
