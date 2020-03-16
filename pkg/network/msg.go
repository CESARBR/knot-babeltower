package network

import "github.com/CESARBR/knot-babeltower/pkg/thing/entities"

// DeviceRegisterRequest represents the incoming register device request message
type DeviceRegisterRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// DeviceRegisteredResponse represents the outgoing register device response message
type DeviceRegisteredResponse struct {
	ID    string  `json:"id"`
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

// SchemaUpdateRequest represents the incoming update schema request message
type SchemaUpdateRequest struct {
	ID     string            `json:"id"`
	Schema []entities.Schema `json:"schema,omitempty"`
}

// SchemaUpdatedResponse represents the outgoing update schema response message
type SchemaUpdatedResponse struct {
	ID string `json:"id"`
}

// DeviceAuthRequest represents the incoming auth device command
type DeviceAuthRequest struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

// DeviceAuthResponse represents the outgoing auth device command response
type DeviceAuthResponse struct {
	ID     string  `json:"id"`
	ErrMsg *string `json:"error"`
}

// DeviceListResponse represents the outgoing list devices command response
type DeviceListResponse struct {
	Things []*entities.Thing `json:"things"`
}

// DataRequest represents the incoming request data command
type DataRequest struct {
	ID        string `json:"id"`
	SensorIds []int  `json:"sensorIds"`
}
