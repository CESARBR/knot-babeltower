package mocks

import (
	"github.com/stretchr/testify/mock"
)

// FakeSessionStore represents a mocking type for session store capabilities.
// It is composed by the testify mock.Mock type to extend the mocking features
// provided by the library.
type FakeSessionStore struct {
	mock.Mock
	GetReturnErr  error
	SaveReturnErr error
	Session       string
}

// Get provides a mock function to get a session avaialble on caching according to the
// user e-mail.
func (fss *FakeSessionStore) Get(email string) (string, error) {
	ret := fss.Called(email)
	return ret.String(0), ret.Error(1)
}

// Save provides a mock function to save a new session to the caching store.
// It receives a key (email) and the respective value (session ID).
func (fss *FakeSessionStore) Save(email, id string) error {
	ret := fss.Called(email, id)
	return ret.Error(0)
}
