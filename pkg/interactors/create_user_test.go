package interactors

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/CESARBR/knot-babeltower/pkg/entities"
)

type FakeCreateUserLogger struct {
}

type FakeUserProxy struct {
	mock.Mock
}

func (fl *FakeCreateUserLogger) Info(...interface{}) {}

func (fl *FakeCreateUserLogger) Infof(string, ...interface{}) {}

func (fl *FakeCreateUserLogger) Debug(...interface{}) {}

func (fl *FakeCreateUserLogger) Warn(...interface{}) {}

func (fl *FakeCreateUserLogger) Error(...interface{}) {}

func (fl *FakeCreateUserLogger) Errorf(string, ...interface{}) {}

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

func (fup *FakeUserProxy) CreateToken(user entities.User) (string, error) {
	return "", nil
}

func TestCreateUser(t *testing.T) {
	testCases := []struct {
		name          string
		email         string
		password      string
		fakeLogger    *FakeCreateUserLogger
		fakeUserProxy *FakeUserProxy
		proxyError    error
		expected      error
	}{
		{
			"shouldCallLogger",
			"fake@email.com",
			"123",
			&FakeCreateUserLogger{},
			&FakeUserProxy{},
			nil,
			nil,
		},
		{
			"shouldRaiseEntitiesExists",
			"fake@email.com",
			"123",
			&FakeCreateUserLogger{},
			&FakeUserProxy{},
			entities.ErrEntityExists{Msg: "User exists"},
			entities.ErrEntityExists{Msg: "mocked msg"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			createUserInteractor := NewCreateUser(tc.fakeLogger, tc.fakeUserProxy)
			user := entities.User{Email: tc.email, Password: tc.password}
			tc.fakeUserProxy.On("Create", user).
				Return(tc.proxyError).Once()

			err := createUserInteractor.Execute(user)
			if err != nil && !assert.IsType(t, err, tc.expected) {
				t.Errorf("Create User failed. Error: %s", err)
				tc.fakeUserProxy.AssertExpectations(t)
				return
			}

			t.Logf("Create user ok")
			tc.fakeUserProxy.AssertExpectations(t)
		})
	}

}
