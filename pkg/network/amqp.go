package network

import "github.com/CESARBR/knot-babeltower/pkg/logging"

// Amqp handles the connection, queues and exchanges declared
type Amqp struct {
	url    string
	logger logging.Logger
}

// NewAmqp constructs the AMQP connection handler
func NewAmqp(url string, logger logging.Logger) *Amqp {
	return &Amqp{url, logger}
}

// Start starts the handler
func (ah *Amqp) Start() {
	ah.logger.Debug("AMQP handler started")
	// TODO: Start amqp connection
}
