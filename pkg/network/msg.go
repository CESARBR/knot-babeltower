package network

import "github.com/CESARBR/knot-babeltower/pkg/thing/entities"

// RegisterRequestMsg is received to register a device
type RegisterRequestMsg struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// UpdateSchemaRequestMsg represents the update schema request message
type UpdateSchemaRequestMsg struct {
	ID     string            `json:"id"`
	Schema []entities.Schema `json:"schema,omitempty"`
}

// RegisterResponseMsg is sent when receive a register request
type RegisterResponseMsg struct {
	ID    string  `json:"id"`
	Token string  `json:"token"`
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
