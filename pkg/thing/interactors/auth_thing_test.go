package interactors

import (
	"testing"

	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
	"github.com/CESARBR/knot-babeltower/pkg/thing/mocks"
	"github.com/stretchr/testify/assert"
)

type AuthThingTestCase struct {
	name           string
	authParam      string
	idParam        string
	expectedErr    error
	fakeLogger     *mocks.FakeLogger
	fakeThingProxy *mocks.FakeThingProxy
}

var atCases = []AuthThingTestCase{
	{
		"authorization token not provided",
		"",
		"8380ba096a091fb9",
		ErrAuthNotProvided,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
	},
	{
		"thing's id not provided",
		"authorization-token",
		"",
		ErrIDNotProvided,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
	},
	{
		"thing 6c0dcd9833b595f9 not found",
		"authorization-token",
		"6c0dcd9833b595f9",
		entities.ErrThingNotFound,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{ReturnErr: entities.ErrThingNotFound},
	},
	{
		"forbidden to authenticate thing",
		"invalid-authorization-token",
		"8380ba096a091fb9",
		entities.ErrThingForbidden,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{ReturnErr: entities.ErrThingForbidden},
	},
	{
		"allowed to authenticate thing",
		"authorization-token",
		"8380ba096a091fb9",
		nil,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{Thing: &entities.Thing{
			ID:    "fc3fcf912d0c290a",
			Token: "token",
			Name:  "thing",
		}},
	},
}

func TestAuthThing(t *testing.T) {
	for _, tc := range atCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.fakeThingProxy.
				On("Get", tc.authParam, tc.idParam).
				Return(tc.fakeThingProxy.Thing, tc.fakeThingProxy.ReturnErr).
				Maybe()

			thingInteractor := NewThingInteractor(tc.fakeLogger, nil, tc.fakeThingProxy, nil)
			err := thingInteractor.Auth(tc.authParam, tc.idParam)

			if tc.authParam == "" {
				msg := tc.expectedErr.Error()
				assert.EqualError(t, err, msg)
			}
			if tc.idParam == "" {
				msg := tc.expectedErr.Error()
				assert.EqualError(t, err, msg)
			}

			tc.fakeThingProxy.AssertExpectations(t)
		})
	}
}
