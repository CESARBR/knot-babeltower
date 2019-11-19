package network

// InMsg represents the message received from the AMQP broker
type InMsg struct {
	Exchange   string
	RoutingKey string
	Headers    map[string]interface{}
	Body       []byte
}

// RegisterRequestMsg is received to register a device
type RegisterRequestMsg struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// RegisterResponseMsg is sent when receive a register request
type RegisterResponseMsg struct {
	ID    string  `json:"id"`
	Token string  `json:"token"`
	Error *string `json:"error"`
}
