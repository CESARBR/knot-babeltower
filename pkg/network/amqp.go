package network

import (
	"gopkg.in/cenkalti/backoff.v3"

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

func (a *Amqp) notifyWhenClosed(started chan bool) {
	errReason := <-a.conn.NotifyClose(make(chan *amqp.Error))
	a.logger.Infof("AMQP connection closed: %s", errReason)
	started <- false
	if errReason != nil {
		err := backoff.Retry(a.connect, backoff.NewExponentialBackOff())
		if err != nil {
			a.logger.Error(err)
			started <- false
			return
		}

		go a.notifyWhenClosed(started)
		started <- true
	}
}

func (a *Amqp) connect() error {
	conn, err := amqp.Dial(a.url)
	if err != nil {
		a.logger.Error(err)
		return err
	}

	a.conn = conn

	channel, err := a.conn.Channel()
	if err != nil {
		a.logger.Error(err)
		return err
	}

	a.logger.Debug("AMQP handler connected")
	a.channel = channel

	return nil
}

// NewAmqp constructs the AMQP connection handler
func NewAmqp(url string, logger logging.Logger) *Amqp {
	return &Amqp{url, logger, nil, nil}
}

// Start starts the handler
func (a *Amqp) Start(started chan bool) {
	err := backoff.Retry(a.connect, backoff.NewExponentialBackOff())
	if err != nil {
		a.logger.Error(err)
		started <- false
		return
	}

	go a.notifyWhenClosed(started)
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
