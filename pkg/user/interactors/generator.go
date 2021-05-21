package interactors

import "github.com/segmentio/ksuid"

// Generator is an interface for generating random identifiers.
type Generator interface {
	ID() string
}

type generator struct{}

// NewGenerator creates a new instance of generator, which implements the Generator interface.
func NewGenerator() Generator {
	return &generator{}
}

// ID creates a new random ID and returns it as a string.
func (r *generator) ID() string {
	randomID := ksuid.New()
	return randomID.String()
}
