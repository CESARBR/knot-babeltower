package mocks

import (
	"github.com/CESARBR/knot-babeltower/pkg/user/entities"
	"github.com/stretchr/testify/mock"
)

// FakeUserProxy represents a mocking type for the user's proxy service
type FakeUserProxy struct {
	mock.Mock
	Token string
	Err   error
}

// Create provides a mock function to create a new user
func (fup *FakeUserProxy) Create(user entities.User) (err error) {
	ret := fup.Called(user)

	rf, ok := ret.Get(0).(func(entities.User) error)
	if ok {
		err = rf(user)
	} else {
		err = ret.Error(0)
	}

	return err
}

// CreateToken provides a mock function to create a new user's token
func (fup *FakeUserProxy) CreateToken(user entities.User) (string, error) {
	args := fup.Called(user)
	return args.String(0), args.Error(1)
}
