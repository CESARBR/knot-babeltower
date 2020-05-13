package interactors

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/CESARBR/knot-babeltower/pkg/mocks"
	"github.com/CESARBR/knot-babeltower/pkg/user/entities"
)

type CreateUserTestCase struct {
	name          string
	email         string
	password      string
	expected      error
	fakeLogger    *mocks.FakeLogger
	fakeUserProxy *mocks.FakeUserProxy
}

var cuCases = []CreateUserTestCase{
	{
		"user successfully created",
		"user@user.com",
		"123456789abcdef",
		nil,
		&mocks.FakeLogger{},
		&mocks.FakeUserProxy{},
	},
	{
		"failed to create user when already exists",
		"fake@email.com",
		"123456789abcdef",
		entities.ErrUserExists,
		&mocks.FakeLogger{},
		&mocks.FakeUserProxy{Err: entities.ErrUserExists},
	},
	{
		"failed to create user when e-mail or password format are invalid",
		"user",
		"123456789abcdef",
		entities.ErrUserBadRequest,
		&mocks.FakeLogger{},
		&mocks.FakeUserProxy{Err: entities.ErrUserBadRequest},
	},
}

func TestCreateUser(t *testing.T) {
	for _, tc := range cuCases {
		t.Run(tc.name, func(t *testing.T) {
			createUserInteractor := NewCreateUser(tc.fakeLogger, tc.fakeUserProxy)
			user := entities.User{Email: tc.email, Password: tc.password}
			tc.fakeUserProxy.On("Create", user).
				Return(tc.fakeUserProxy.Err).Once()

			err := createUserInteractor.Execute(user)
			if err != nil && !assert.IsType(t, err, tc.expected) {
				t.Errorf("create user failed. Error: %s", err)
				tc.fakeUserProxy.AssertExpectations(t)
				return
			}

			t.Logf("create user ok")
			tc.fakeUserProxy.AssertExpectations(t)
		})
	}

}
