package entities

// Thing represents the thing domain entity
type Thing struct {
	ID     string   `json:"id"`
	Token  string   `json:"token,omitempty"`
	Name   string   `json:"name,omitempty"`
	Config []Config `json:"config,omitempty"`
}
