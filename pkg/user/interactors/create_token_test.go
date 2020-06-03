package interactors

import (
	"testing"

	"github.com/CESARBR/knot-babeltower/pkg/mocks"
	"github.com/CESARBR/knot-babeltower/pkg/user/entities"
	"github.com/stretchr/testify/assert"
)

type createTokenResponse struct {
	token string
	err   error
}

type CreateTokenTestCase struct {
	name           string
	email          string
	password       string
	expected       createTokenResponse
	fakeLogger     *mocks.FakeLogger
	fakeUsersProxy *mocks.FakeUsersProxy
}

var ctCases = []CreateTokenTestCase{
	{
		"user token successfully created",
		"user@user.com",
		"123456789abcdef",
		createTokenResponse{"mocked-token", nil},
		&mocks.FakeLogger{},
		&mocks.FakeUsersProxy{Token: "mocked-token"},
	},
	{
		"failed to create user token when e-mail or password format are invalid",
		"user",
		"123456789abcdef",
		createTokenResponse{"", entities.ErrUserBadRequest},
		&mocks.FakeLogger{},
		&mocks.FakeUsersProxy{Err: entities.ErrUserBadRequest},
	},
	{
		"failed to create user token when unauthorized",
		"user@user.com",
		"abcdef",
		createTokenResponse{"", entities.ErrUserForbidden},
		&mocks.FakeLogger{},
		&mocks.FakeUsersProxy{Err: entities.ErrUserForbidden},
	},
}

func TestCreateToken(t *testing.T) {
	for _, tc := range ctCases {
		t.Run(tc.name, func(t *testing.T) {
			user := entities.User{Email: tc.email, Password: tc.password}
			tc.fakeUsersProxy.
				On("CreateToken", user).
				Return(tc.fakeUsersProxy.Token, tc.fakeUsersProxy.Err)

			createTokenInteractor := NewCreateToken(tc.fakeLogger, tc.fakeUsersProxy)
			token, err := createTokenInteractor.Execute(user)
			if err != nil {
				assert.Equal(t, tc.expected.token, token)
				return
			}

			assert.Nil(t, err)
		})
	}
}
