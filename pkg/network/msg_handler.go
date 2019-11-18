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
	logger       logging.Logger
	amqp         *Amqp
	msgPublisher *MsgPublisher
}

func (mc *MsgHandler) handleRegisterMsg(body []byte) error {
	msgParsed := RegisterRequestMsg{}
	err := json.Unmarshal(body, &msgParsed)
	if err != nil {
		return err
	}

	response := RegisterResponseMsg{ID: msgParsed.ID, Token: "secret", Error: nil}
	return mc.msgPublisher.SendRegisterDevice(response)
}

func (mc *MsgHandler) onMsgReceived(msgChan chan InMsg) {
	for {
		msg := <-msgChan
		mc.logger.Infof("Exchange: %s, routing key: %s", msg.Exchange, msg.RoutingKey)
		mc.logger.Infof("Message received: %s", string(msg.Body))

		switch msg.RoutingKey {
		case "device.register":
			err := mc.handleRegisterMsg(msg.Body)
			if err != nil {
				mc.logger.Error(err)
				continue
			}
		}
	}
}

// NewMsgHandler constructs the MsgHandler
func NewMsgHandler(logger logging.Logger, amqp *Amqp, msgPublisher *MsgPublisher) *MsgHandler {
	return &MsgHandler{logger, amqp, msgPublisher}
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
