package interactors

import (
	"errors"
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
	fakePublisher  *mocks.FakePublisher
	fakeConnector  *mocks.FakeConnector
}

var errMsg1 = "Thing 6c0dcd9833b595f9 not found"
var errMsg2 = "Forbidden to authenticate thing"

var atCases = []AuthThingTestCase{
	{
		"authorization key not provided",
		"",
		"8380ba096a091fb9",
		errors.New("authorization key not provided"),
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
		&mocks.FakeConnector{},
	},
	{
		"thing's id not provided",
		"authorization-token",
		"",
		errors.New("thing's id not provided"),
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
		&mocks.FakeConnector{},
	},
	{
		"thing 6c0dcd9833b595f9 not found",
		"authorization-token",
		"6c0dcd9833b595f9",
		entities.ErrThingNotFound{ID: "6c0dcd9833b595f9"},
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{ReturnErr: entities.ErrThingNotFound{ID: "6c0dcd9833b595f9"}},
		&mocks.FakePublisher{ErrMsg: &errMsg1},
		&mocks.FakeConnector{},
	},
	{
		"forbidden to authenticate the thing",
		"invalid-authorization-token",
		"8380ba096a091fb9",
		&entities.ErrThingForbidden{},
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{ReturnErr: &entities.ErrThingForbidden{}},
		&mocks.FakePublisher{ErrMsg: &errMsg2},
		&mocks.FakeConnector{},
	},
	{
		"allowed to authenticate the thing",
		"authorization-token",
		"8380ba096a091fb9",
		nil,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{Thing: &entities.Thing{
			ID:    "fc3fcf912d0c290a",
			Token: "token",
			Name:  "thing",
		}},
		&mocks.FakePublisher{},
		&mocks.FakeConnector{},
	},
	{
		"failed to send the authenticate response",
		"authorization-token",
		"8380ba096a091fb9",
		errors.New("failed to send authenticate response"),
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{Thing: &entities.Thing{
			ID:    "fc3fcf912d0c290a",
			Token: "token",
			Name:  "thing",
		}},
		&mocks.FakePublisher{SendError: errors.New("failed to send authenticate response")},
		&mocks.FakeConnector{},
	},
	{
		"authenticate response successfully sent",
		"authorization-token",
		"8380ba096a091fb9",
		nil,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{Thing: &entities.Thing{
			ID:    "fc3fcf912d0c290a",
			Token: "token",
			Name:  "thing",
		}},
		&mocks.FakePublisher{},
		&mocks.FakeConnector{},
	},
}

func TestAuthThing(t *testing.T) {
	for _, tc := range atCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.fakeThingProxy.
				On("Get", tc.authParam, tc.idParam).
				Return(tc.fakeThingProxy.Thing, tc.fakeThingProxy.ReturnErr).
				Maybe()
			tc.fakePublisher.
				On("SendAuthStatus", tc.idParam, tc.fakePublisher.ErrMsg).
				Return(tc.fakeConnector.SendError).
				Maybe()

			thingInteractor := NewThingInteractor(tc.fakeLogger, tc.fakePublisher, tc.fakeThingProxy, tc.fakeConnector)
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
			tc.fakePublisher.AssertExpectations(t)
		})
	}
}
