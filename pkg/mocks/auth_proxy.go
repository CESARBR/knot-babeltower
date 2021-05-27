package mocks

import (
	"github.com/CESARBR/knot-babeltower/pkg/user/entities"
	"github.com/stretchr/testify/mock"
)

// FakeAuthProxy represents a mocking type for the user's proxy service
type FakeAuthProxy struct {
	mock.Mock
	Token string
	Err   error
}

// CreateAppToken provides a mock function to create a new application token
func (fup *FakeAuthProxy) CreateAppToken(user entities.User, duration int) (string, error) {
	args := fup.Called(user, duration)
	return args.String(0), args.Error(1)
}
