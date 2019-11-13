package network

// InMsg represents the message received from the AMQP broker
type InMsg struct {
	Exchange   string
	RoutingKey string
	Body       []byte
}
