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
	duration       int
	expected       createTokenResponse
	fakeLogger     *mocks.FakeLogger
	fakeUsersProxy *mocks.FakeUsersProxy
	fakeAuthProxy  *mocks.FakeAuthProxy
}

var (
	mockedUser = entities.User{
		Email:    "user@knot.com",
		Password: "123qwe123qwe",
		Token:    "user-access-token",
	}
	mockedToken   = "mocked-token"
	errUsersProxy = errors.New("fail to create a user token")
	errAuthProxy  = errors.New("fail to create a app token")
)

var ctCases = []createTokenTestCase{
	{
		"user token successfully created",
		mockedUser,
		"user",
		0,
		createTokenResponse{mockedToken, nil},
		&mocks.FakeLogger{},
		&mocks.FakeUsersProxy{Token: mockedToken},
		&mocks.FakeAuthProxy{},
	},
	{
		"app token successfully created",
		mockedUser,
		"app",
		3600,
		createTokenResponse{mockedToken, nil},
		&mocks.FakeLogger{},
		&mocks.FakeUsersProxy{},
		&mocks.FakeAuthProxy{Token: mockedToken},
	},
	{
		"failed to create user token when something goes wrong in usersProxy",
		mockedUser,
		"user",
		0,
		createTokenResponse{"", errUsersProxy},
		&mocks.FakeLogger{},
		&mocks.FakeUsersProxy{Err: errUsersProxy},
		&mocks.FakeAuthProxy{},
	},
	{
		"failed to create app token when something goes wrong in authProxy",
		mockedUser,
		"app",
		3600,
		createTokenResponse{"", errAuthProxy},
		&mocks.FakeLogger{},
		&mocks.FakeUsersProxy{},
		&mocks.FakeAuthProxy{Err: errAuthProxy},
	},
	{
		"fail to create a token when the tokenType is invalid",
		mockedUser,
		"invalid-token-type",
		3600,
		createTokenResponse{"", entities.ErrInvalidTokenType},
		&mocks.FakeLogger{},
		&mocks.FakeUsersProxy{},
		&mocks.FakeAuthProxy{},
	},
}

func TestCreateToken(t *testing.T) {
	for _, tc := range ctCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.fakeUsersProxy.
				On("CreateToken", tc.user).
				Return(tc.fakeUsersProxy.Token, tc.fakeUsersProxy.Err)
			tc.fakeAuthProxy.
				On("CreateAppToken", tc.user, tc.duration).
				Return(tc.fakeAuthProxy.Token, tc.fakeUsersProxy.Err)

			createTokenInteractor := NewCreateToken(tc.fakeLogger, tc.fakeUsersProxy, tc.fakeAuthProxy)
			token, err := createTokenInteractor.Execute(tc.user, tc.tokenType, tc.duration)

			assert.Equal(t, tc.expected.token, token)
			if err != nil {
				assert.Equal(t, errors.Is(err, tc.expected.err), true)
			}

			if tc.tokenType == "user" {
				tc.fakeUsersProxy.AssertExpectations(t)
			}
			if tc.tokenType == "app" {
				tc.fakeAuthProxy.AssertExpectations(t)
			}
		})
	}
}
