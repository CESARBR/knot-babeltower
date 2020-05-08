package network

import (
	"fmt"

	"github.com/cenkalti/backoff/v4"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/streadway/amqp"
)

// QueueService is the interface that handles remote queue service
type QueueService interface {
	Start(started chan bool)
	Stop()
	GetSender() AmqpSender
	GetReceiver() AmqpReceiver
}

// AmqpSender is the interface to publish amqp messages
type AmqpSender interface {
	PublishPersistentMessage(exchange, exchangeType, key string, msg MessageSerializer, options *MessageOptions) error
}

// AmqpReceiver is the interface to receive amqp messages
type AmqpReceiver interface {
	OnMessage(msgChan chan InMsg, queueName, exchangeName, exchangeType, key string) error
}

// Amqp handles the connection, queues and exchanges declared
type Amqp struct {
	url     string
	logger  logging.Logger
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   *amqp.Queue
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
}

// NewAmqp constructs the AMQP connection handler
func NewAmqp(url string, logger logging.Logger) *Amqp {
	return &Amqp{url, logger, nil, nil, nil}
}

// GetSender returns the sender
func (a *Amqp) GetSender() AmqpSender {
	return a
}

// GetReceiver returns the receiver
func (a *Amqp) GetReceiver() AmqpReceiver {
	return a
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

// PublishPersistentMessage sends a persistent message to RabbitMQ
func (a *Amqp) PublishPersistentMessage(exchange, exchangeType, key string, msg MessageSerializer, options *MessageOptions) error {
	var headers map[string]interface{}
	var corrID string

	if options != nil {
		headers = map[string]interface{}{
			"Authorization": options.Authorization,
		}
		corrID = options.CorrelationID
	}

	body, err := msg.Serialize()
	if err != nil {
		return fmt.Errorf("error serializing message: %w", err)
	}

	err = a.declareExchange(exchange, exchangeType)
	if err != nil {
		return fmt.Errorf("error declaring exchange: %w", err)
	}

	err = a.channel.Publish(
		exchange,
		key,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			Headers:         headers,
			ContentType:     "text/plain",
			ContentEncoding: "",
			DeliveryMode:    amqp.Persistent,
			Priority:        0,
			CorrelationId:   corrID,
			Body:            body,
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
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		a.logger.Error(err)
		return err
	}

	deliveries, err := a.channel.Consume(
		queueName,
		"",    // consumerTag
		true,  // noAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // arguments
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
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	)
}

func (a *Amqp) declareQueue(name string) error {
	queue, err := a.channel.QueueDeclare(
		name,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)

	a.queue = &queue
	return err
}

func convertDeliveryToInMsg(deliveries <-chan amqp.Delivery, outMsg chan InMsg) {
	for d := range deliveries {
		outMsg <- InMsg{d.Exchange, d.RoutingKey, d.ReplyTo, d.CorrelationId, d.Headers, d.Body}
	}
}
