package interactors

import (
	"errors"
	"testing"

	"github.com/CESARBR/knot-babeltower/pkg/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type FakeCreateTokenLogger struct{}

type FakeTokenProxy struct {
	mock.Mock
}

type CreateUserResponse struct {
	token string
	err   error
}

type CreateTokenTestCase struct {
	name           string
	email          string
	password       string
	fakeLogger     *FakeCreateTokenLogger
	fakeTokenProxy *FakeTokenProxy
	expected       CreateUserResponse
}

func (fl *FakeCreateTokenLogger) Info(...interface{}) {}

func (fl *FakeCreateTokenLogger) Infof(string, ...interface{}) {}

func (fl *FakeCreateTokenLogger) Debug(...interface{}) {}

func (fl *FakeCreateTokenLogger) Warn(...interface{}) {}

func (fl *FakeCreateTokenLogger) Error(...interface{}) {}

func (fl *FakeCreateTokenLogger) Errorf(string, ...interface{}) {}

func (ftp *FakeTokenProxy) Create(user entities.User) error {
	return nil
}

func (ftp *FakeTokenProxy) CreateToken(user entities.User) (token string, err error) {
	args := ftp.Called(user)
	return args.String(0), args.Error(1)
}

var cases = []CreateTokenTestCase{
	{
		"success",
		"fake@email.com",
		"123456",
		&FakeCreateTokenLogger{},
		&FakeTokenProxy{},
		CreateUserResponse{"mocked-token", nil},
	},
	{
		"create-token-failed",
		"fake@email.com",
		"123456",
		&FakeCreateTokenLogger{},
		&FakeTokenProxy{},
		CreateUserResponse{"", errors.New("failed to create token")},
	},
}

func TestCreateToken(t *testing.T) {
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			user := entities.User{Email: tc.email, Password: tc.password}
			tc.fakeTokenProxy.
				On("CreateToken", user).
				Return(tc.expected.token, tc.expected.err)

			createTokenInteractor := NewCreateToken(tc.fakeLogger, tc.fakeTokenProxy)
			token, err := createTokenInteractor.Execute(user)
			if err != nil {
				assert.Equal(t, tc.expected.token, token)
				return
			}

			assert.Nil(t, err)
		})
	}
}
