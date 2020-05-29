package mocks

import (
	"github.com/CESARBR/knot-babeltower/pkg/user/entities"
	"github.com/stretchr/testify/mock"
)

// FakeAuthnProxy represents a mocking type for the user's proxy service
type FakeAuthnProxy struct {
	mock.Mock
	Token string
	Err   error
}

// CreateAppToken provides a mock function to create a new application token
func (fup *FakeAuthnProxy) CreateAppToken(user entities.User) (string, error) {
	args := fup.Called(user)
	return args.String(0), args.Error(1)
}
