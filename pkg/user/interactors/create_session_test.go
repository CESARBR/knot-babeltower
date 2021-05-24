package interactors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/CESARBR/knot-babeltower/pkg/jwt"
	"github.com/CESARBR/knot-babeltower/pkg/mocks"
	"github.com/CESARBR/knot-babeltower/pkg/user/entities"
	"github.com/stretchr/testify/assert"
)

type expectedReturn struct {
	id  string
	err error
}

type createSession struct {
	name             string
	authorization    string
	email            string
	expected         expectedReturn
	fakeGenerator    *mocks.FakeGenerator
	fakeSessionStore *mocks.FakeSessionStore
	fakeThingProxy   *mocks.FakeThingProxy
}

var (
	validToken          = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjE1NTI2NDUsImlhdCI6MTYyMTUxNjY0NSwiaXNzIjoibWFpbmZsdXguYXV0aG4iLCJzdWIiOiJ1c2VyMUBjZXNhci5vcmcuYnIiLCJ0eXBlIjowfQ.ndhiTDqmB0sFyuiave2cyvf26ayHUDSO52yLAvmOn_Q"
	errGetSession       = errors.New("failed to get user session")
	errTokenValidation  = errors.New("failed to authenticate token")
	errIDGeneration     = errors.New("failed to generate ID")
	errSaveSessionStore = errors.New("failed to save session on store")
)

var csCases = []createSession{
	{
		"new user session successfully created",
		validToken,
		"user1@cesar.org.br",
		expectedReturn{"session-id", nil},
		&mocks.FakeGenerator{ReturnID: "session-id"},
		&mocks.FakeSessionStore{},
		&mocks.FakeThingProxy{},
	},
	{
		"return existing session successfully",
		validToken,
		"user1@cesar.org.br",
		expectedReturn{"session-id", nil},
		&mocks.FakeGenerator{},
		&mocks.FakeSessionStore{Session: "session-id"},
		&mocks.FakeThingProxy{},
	},
	{
		"failed to create session if token is invalid",
		validToken,
		"user1@cesar.org.br",
		expectedReturn{"", entities.ErrTokenForbidden},
		&mocks.FakeGenerator{},
		&mocks.FakeSessionStore{},
		&mocks.FakeThingProxy{ReturnErr: errTokenValidation},
	},
	{
		"failed to create session if token can't be parsed",
		"invalid-token",
		"user1@cesar.org.br",
		expectedReturn{"", jwt.ErrParseToken},
		&mocks.FakeGenerator{},
		&mocks.FakeSessionStore{},
		&mocks.FakeThingProxy{},
	},
	{
		"failed to get session if it exists",
		validToken,
		"user1@cesar.org.br",
		expectedReturn{"", errGetSession},
		&mocks.FakeGenerator{},
		&mocks.FakeSessionStore{GetReturnErr: errGetSession},
		&mocks.FakeThingProxy{},
	},
	{
		"failed generate session ID",
		validToken,
		"user1@cesar.org.br",
		expectedReturn{"", errIDGeneration},
		&mocks.FakeGenerator{ReturnErr: errIDGeneration},
		&mocks.FakeSessionStore{},
		&mocks.FakeThingProxy{},
	},
	{
		"failed to save a new session to session store",
		validToken,
		"user1@cesar.org.br",
		expectedReturn{"", errSaveSessionStore},
		&mocks.FakeGenerator{},
		&mocks.FakeSessionStore{SaveReturnErr: errSaveSessionStore},
		&mocks.FakeThingProxy{},
	},
}

func TestCreateSsesion(t *testing.T) {
	for _, tc := range csCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.fakeThingProxy.
				On("List", tc.authorization).
				Return(tc.fakeThingProxy.Things, tc.fakeThingProxy.ReturnErr).
				Maybe()
			tc.fakeSessionStore.
				On("Get", tc.email).
				Return(tc.fakeSessionStore.Session, tc.fakeSessionStore.GetReturnErr).
				Maybe()
			tc.fakeGenerator.
				On("ID").
				Return(tc.fakeGenerator.ReturnID).
				Maybe()
			tc.fakeSessionStore.
				On("Save", tc.email, tc.expected.id).
				Return(tc.fakeSessionStore.SaveReturnErr).
				Maybe()

			createSessionInteractor := NewCreateSession(tc.fakeThingProxy, tc.fakeGenerator, tc.fakeSessionStore)
			id, err := createSessionInteractor.Execute(tc.authorization)

			assert.Equal(t, tc.expected.id, id)
			if err != nil {
				assert.Equal(t, errors.Is(err, tc.expected.err), true)
				fmt.Println(err)
			}

			tc.fakeThingProxy.AssertExpectations(t)
			tc.fakeGenerator.AssertExpectations(t)
			tc.fakeSessionStore.AssertExpectations(t)
		})
	}
}
