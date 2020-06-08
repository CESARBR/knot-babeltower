package interactors

import (
	"errors"
	"testing"

	"github.com/CESARBR/knot-babeltower/pkg/mocks"
	"github.com/CESARBR/knot-babeltower/pkg/user/entities"
	"github.com/stretchr/testify/assert"
)

type createTokenResponse struct {
	token string
	err   error
}

type createTokenTestCase struct {
	name           string
	user           entities.User
	tokenType      string
	expected       createTokenResponse
	fakeLogger     *mocks.FakeLogger
	fakeUsersProxy *mocks.FakeUsersProxy
	fakeAuthnProxy *mocks.FakeAuthnProxy
}

var (
	mockedUser = entities.User{
		Email:    "user@knot.com",
		Password: "123qwe123qwe",
		Token:    "user-access-token",
	}
	mockedToken   = "mocked-token"
	errUsersProxy = errors.New("fail to create a user token")
	errAuthnProxy = errors.New("fail to create a app token")
)

var ctCases = []createTokenTestCase{
	{
		"user token successfully created",
		mockedUser,
		"user",
		createTokenResponse{mockedToken, nil},
		&mocks.FakeLogger{},
		&mocks.FakeUsersProxy{Token: mockedToken},
		&mocks.FakeAuthnProxy{},
	},
	{
		"app token successfully created",
		mockedUser,
		"app",
		createTokenResponse{mockedToken, nil},
		&mocks.FakeLogger{},
		&mocks.FakeUsersProxy{},
		&mocks.FakeAuthnProxy{Token: mockedToken},
	},
	{
		"failed to create user token when something goes wrong in usersProxy",
		mockedUser,
		"user",
		createTokenResponse{"", errUsersProxy},
		&mocks.FakeLogger{},
		&mocks.FakeUsersProxy{Err: errUsersProxy},
		&mocks.FakeAuthnProxy{},
	},
	{
		"failed to create app token when something goes wrong in authnProxy",
		mockedUser,
		"app",
		createTokenResponse{"", errAuthnProxy},
		&mocks.FakeLogger{},
		&mocks.FakeUsersProxy{},
		&mocks.FakeAuthnProxy{Err: errAuthnProxy},
	},
	{
		"fail to create a token when the tokenType is invalid",
		mockedUser,
		"invalid-token-type",
		createTokenResponse{"", entities.ErrInvalidTokenType},
		&mocks.FakeLogger{},
		&mocks.FakeUsersProxy{},
		&mocks.FakeAuthnProxy{},
	},
}

func TestCreateToken(t *testing.T) {
	for _, tc := range ctCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.fakeUsersProxy.
				On("CreateToken", tc.user).
				Return(tc.fakeUsersProxy.Token, tc.fakeUsersProxy.Err)
			tc.fakeAuthnProxy.
				On("CreateAppToken", tc.user).
				Return(tc.fakeAuthnProxy.Token, tc.fakeUsersProxy.Err)

			createTokenInteractor := NewCreateToken(tc.fakeLogger, tc.fakeUsersProxy, tc.fakeAuthnProxy)
			token, err := createTokenInteractor.Execute(tc.user, tc.tokenType)

			assert.Equal(t, tc.expected.token, token)
			if err != nil {
				assert.Equal(t, errors.Is(err, tc.expected.err), true)
			}

			if tc.tokenType == "user" {
				tc.fakeUsersProxy.AssertExpectations(t)
			}
			if tc.tokenType == "app" {
				tc.fakeAuthnProxy.AssertExpectations(t)
			}
		})
	}
}
