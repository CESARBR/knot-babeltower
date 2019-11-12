package network

import (
	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/streadway/amqp"
)

// Amqp handles the connection, queues and exchanges declared
type Amqp struct {
	url     string
	logger  logging.Logger
	conn    *amqp.Connection
	channel *amqp.Channel
}

func (a *Amqp) notifyWhenClosed() {
	errReason := <-a.conn.NotifyClose(make(chan *amqp.Error))
	a.logger.Infof("AMQP connection closed: %s", errReason)
	// TODO: try to reconnect
}

// NewAmqp constructs the AMQP connection handler
func NewAmqp(url string, logger logging.Logger) *Amqp {
	return &Amqp{url, logger, nil, nil}
}

// Start starts the handler
func (a *Amqp) Start(started chan bool) {
	conn, err := amqp.Dial(a.url)
	if err != nil {
		// TODO: try to reconnect
		a.logger.Error(err)
		started <- false
		return
	}

	a.conn = conn
	go a.notifyWhenClosed()

	channel, err := conn.Channel()
	if err != nil {
		// TODO: try to create channel again
		a.logger.Error(err)
		started <- false
		return
	}

	a.logger.Debug("AMQP handler connected")
	a.channel = channel
	started <- true
}

// Stop closes the connection started
func (a *Amqp) Stop() {
	if a.conn != nil && !a.conn.IsClosed() {
		a.conn.Close()
	}

	if a.channel != nil {
		a.channel.Close()
	}

	a.logger.Debug("AMQP handler stopped")
}
