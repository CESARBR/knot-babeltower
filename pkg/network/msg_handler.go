package network

import "github.com/CESARBR/knot-babeltower/pkg/logging"

// MsgHandler handle messages received from a service
type MsgHandler struct {
	logger logging.Logger
	amqp   *Amqp
}

// NewMsgHandler constructs the MsgHandler
func NewMsgHandler(logger logging.Logger, amqp *Amqp) *MsgHandler {
	return &MsgHandler{logger, amqp}
}

// Start starts to listen for messages
func (mc *MsgHandler) Start(started chan bool) {
	mc.logger.Debug("Msg handler started")
	started <- true
}

// Stop stops to listen for messages
func (mc *MsgHandler) Stop() {
	mc.logger.Debug("Msg handler stopped")
}
