package network

import (
	"encoding/json"
	"fmt"

	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
)

// MessageSerializer represents a interface for KNoT messages
type MessageSerializer interface {
	Serialize() ([]byte, error)
}

// message represents a KNoT message
type message struct {
	Payload interface{}
}

// DeviceRegisterRequest represents the incoming register device request message
type DeviceRegisterRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// DeviceRegisteredResponse represents the outgoing register device response message
type DeviceRegisteredResponse struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Token string  `json:"token"`
	Error *string `json:"error"`
}

// DeviceUnregisterRequest represents the incoming unregister device request message
type DeviceUnregisterRequest struct {
	ID string `json:"id"`
}

// DeviceUnregisteredResponse represents the outgoing unregister device response message
type DeviceUnregisteredResponse struct {
	ID    string  `json:"id"`
	Error *string `json:"error"`
}

// ConfigUpdateRequest represents the incoming update config request message
type ConfigUpdateRequest struct {
	ID     string            `json:"id"`
	Config []entities.Config `json:"config,omitempty"`
}

// ConfigUpdatedResponse represents the outgoing update config response message
type ConfigUpdatedResponse struct {
	ID      string            `json:"id"`
	Config  []entities.Config `json:"config,omitempty"`
	Changed bool              `json:"changed"`
	Error   *string           `json:"error"`
}

// DeviceAuthRequest represents the incoming auth device command
type DeviceAuthRequest struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

// DeviceAuthResponse represents the outgoing auth device command response
type DeviceAuthResponse struct {
	ID    string  `json:"id"`
	Error *string `json:"error"`
}

// DeviceListResponse represents the outgoing list devices command response
type DeviceListResponse struct {
	Things []*entities.Thing `json:"devices"`
	Error  *string           `json:"error"`
}

// DataRequest represents the incoming request data command
type DataRequest struct {
	ID        string `json:"id"`
	SensorIds []int  `json:"sensorIds"`
}

// DataUpdate represents the incoming update data command
type DataUpdate struct {
	ID   string          `json:"id"`
	Data []entities.Data `json:"data"`
}

// DataSent represents the data received from the things
type DataSent struct {
	ID   string          `json:"id"`
	Data []entities.Data `json:"data"`
}

// CreateSessionRequest represents session creation request
type CreateSessionRequest struct {
	Token string `json:"token"`
}

// CreateSessionResponse represents session creation response
type CreateSessionResponse struct {
	ID string `json:"id"`
}

// NewMessage creates a message
func NewMessage(msg interface{}) MessageSerializer {
	return message{Payload: msg}
}

// Serialize serializes the message in a byte stream
func (m message) Serialize() ([]byte, error) {
	data, err := json.Marshal(&m.Payload)
	if err != nil {
		return nil, fmt.Errorf("error econding JSON message: %w", err)
	}
	return data, nil
}
