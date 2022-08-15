package network

import (
	"fmt"

	"github.com/cenkalti/backoff/v4"

	"sync"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/streadway/amqp"
)

var exchangeLock *sync.Mutex = &sync.Mutex{}

const (
	mandatory          = false
	immediate          = false
	noWait             = false
	durable            = true
	exclusive          = false
	noAck              = true
	deleteWhenUnused   = false
	deleteWhenComplete = false
	noLocal            = false
)

// Amqp handles the connection, queues and exchanges declared
type Amqp struct {
	url               string
	logger            logging.Logger
	conn              *amqp.Connection
	channel           *amqp.Channel
	queue             *amqp.Queue
	declaredExchanges map[string]struct{}
}

// InMsg represents the message received from the AMQP broker
type InMsg struct {
	Exchange      string
	RoutingKey    string
	ReplyTo       string
	CorrelationID string
	Headers       map[string]interface{}
	Body          []byte
}

// MessageOptions represents the message publishing options
type MessageOptions struct {
	Authorization string
	CorrelationID string
	Expiration    string
}

// NewAmqp constructs the AMQP connection handler
func NewAmqp(url string, logger logging.Logger) *Amqp {
	declaredExchanges := make(map[string]struct{})
	return &Amqp{url, logger, nil, nil, nil, declaredExchanges}
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
		defer a.conn.Close()
	}

	if a.channel != nil {
		defer a.channel.Close()
	}

	a.logger.Debug("AMQP handler stopped")
}

// PublishPersistentMessage sends a persistent message to RabbitMQ
func (a *Amqp) PublishPersistentMessage(exchange, exchangeType, key string, msg MessageSerializer, options *MessageOptions) error {
	var headers map[string]interface{}
	var corrID, expTime string

	if options != nil {
		headers = map[string]interface{}{
			"Authorization": options.Authorization,
		}
		corrID = options.CorrelationID
		expTime = options.Expiration
	}

	body, err := msg.Serialize()
	if err != nil {
		return fmt.Errorf("error serializing message: %w", err)
	}

	if !a.exchangeAlreadyDeclared(exchange) {
		err = a.declareExchange(exchange, exchangeType)
		if err != nil {
			return fmt.Errorf("error declaring exchange: %w", err)
		} else {
			exchangeLock.Lock()
			a.declaredExchanges[exchange] = struct{}{}
			exchangeLock.Unlock()
		}
	}
	err = a.channel.Publish(
		exchange,
		key,
		mandatory,
		immediate,
		amqp.Publishing{
			Headers:         headers,
			ContentType:     "text/plain",
			ContentEncoding: "",
			DeliveryMode:    amqp.Persistent,
			Priority:        0,
			CorrelationId:   corrID,
			Body:            body,
			Expiration:      expTime,
		},
	)
	if err != nil {
		return fmt.Errorf("error publishing message in channel: %w", err)
	}

	return nil
}

// OnMessage receive messages and put them on channel
func (a *Amqp) OnMessage(msgChan chan InMsg, queueName, exchangeName, exchangeType, key string) error {
	err := a.declareExchange(exchangeName, exchangeType)
	if err != nil {
		a.logger.Error(err)
		return err
	}

	err = a.declareQueue(queueName)
	if err != nil {
		a.logger.Error(err)
		return err
	}

	err = a.channel.QueueBind(
		queueName,
		key,
		exchangeName,
		noWait,
		nil, // arguments
	)
	if err != nil {
		a.logger.Error(err)
		return err
	}

	deliveries, err := a.channel.Consume(
		queueName,
		"", // consumerTag
		noAck,
		exclusive,
		noLocal,
		noWait,
		nil, // arguments
	)
	if err != nil {
		a.logger.Error(err)
		return err
	}

	go convertDeliveryToInMsg(deliveries, msgChan)

	return nil
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

func (a *Amqp) declareExchange(name, exchangeType string) error {
	return a.channel.ExchangeDeclare(
		name,
		exchangeType,
		durable,
		deleteWhenComplete,
		false, // internal
		noWait,
		nil, // arguments
	)
}

func (a *Amqp) declareQueue(name string) error {
	queue, err := a.channel.QueueDeclare(
		name,
		durable,
		deleteWhenUnused,
		exclusive,
		noWait,
		nil, // arguments
	)

	a.queue = &queue
	return err
}

func (a *Amqp) exchangeAlreadyDeclared(exchangeName string) bool {
	exchangeLock.Lock()
	_, ok := a.declaredExchanges[exchangeName]
	exchangeLock.Unlock()
	return ok
}

func convertDeliveryToInMsg(deliveries <-chan amqp.Delivery, outMsg chan InMsg) {
	for d := range deliveries {
		outMsg <- InMsg{d.Exchange, d.RoutingKey, d.ReplyTo, d.CorrelationId, d.Headers, d.Body}
	}
}
