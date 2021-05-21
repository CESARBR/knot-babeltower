package mocks

import (
	"github.com/stretchr/testify/mock"
)

// FakeGenerator represents a mocking type for the ID generator
type FakeGenerator struct {
	mock.Mock
	ReturnErr error
	ReturnID  string
}

// Create provides a mock function to create a random ID
func (fg *FakeGenerator) ID() string {
	ret := fg.Called()
	return ret.String(0)
}
