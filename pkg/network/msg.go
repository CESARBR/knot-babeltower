package network

import "github.com/CESARBR/knot-babeltower/pkg/thing/entities"

// RegisterRequestMsg represents the incoming register device request message
type RegisterRequestMsg struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// UnregisterRequestMsg represents the incoming unregister device request message
type UnregisterRequestMsg struct {
	ID string `json:"id"`
}

// UpdateSchemaRequestMsg represents the update schema request message
type UpdateSchemaRequestMsg struct {
	ID     string            `json:"id"`
	Schema []entities.Schema `json:"schema,omitempty"`
}

// RegisterResponseMsg represents the outgoing register device response message
type RegisterResponseMsg struct {
	ID    string  `json:"id"`
	Token string  `json:"token"`
	Error *string `json:"error"`
}

// UnregisterResponseMsg represents the outgoing unregister device response message
type UnregisterResponseMsg struct {
	ID    string  `json:"id"`
	Error *string `json:"error"`
}

// UpdateSchemaRequest represents the update schema request message
type UpdateSchemaRequest struct {
	ID     string            `json:"id"`
	Schema []entities.Schema `json:"schema,omitempty"`
}

// UpdatedSchemaResponse represents the update schema response message
type UpdatedSchemaResponse struct {
	ID string `json:"id"`
}

// ListThingsResponse represents the list things response
type ListThingsResponse struct {
	Things []*entities.Thing `json:"things"`
}

// RequestDataCommand represents the request data command
type RequestDataCommand struct {
	ID        string `json:"id"`
	SensorIds []int  `json:"sensorIds"`
}

// AuthThingCommand represents the auth device command
type AuthThingCommand struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

// AuthThingResponse represents the auth device command response
type AuthThingResponse struct {
	ID     string  `json:"id"`
	ErrMsg *string `json:"error"`
}
