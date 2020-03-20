package interactors

import (
	"errors"
	"testing"

	"github.com/CESARBR/knot-babeltower/pkg/thing/mocks"
	"github.com/stretchr/testify/assert"
)

type UnregisterThingTestCase struct {
	name           string
	authParam      string
	idParam        string
	fakeLogger     *mocks.FakeLogger
	fakeThingProxy *mocks.FakeThingProxy
	fakePublisher  *mocks.FakePublisher
	fakeConnector  *mocks.FakeConnector
	expectedErrMsg string
}

var unregisterAtCases = []UnregisterThingTestCase{
	{
		"authorization key not provided",
		"",
		"thing-id",
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
		&mocks.FakeConnector{},
		"authorization key not provided",
	},
	{
		"thing's id not provided",
		"authorization-token",
		"",
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
		&mocks.FakeConnector{},
		"thing's id not provided",
	},
	{
		"thing's id not found",
		"authorization-token",
		"thing-id",
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{ReturnErr: errors.New("thing's id not found")},
		&mocks.FakePublisher{Err: errors.New("thing's id not found")},
		&mocks.FakeConnector{},
		"thing's id not found",
	},
	{
		"unable to unregister thing",
		"authorization-token",
		"thing-id",
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{ReturnErr: errors.New("unable to unregister thing")},
		&mocks.FakePublisher{Err: errors.New("unable to unregister thing")},
		&mocks.FakeConnector{},
		"unable to unregister thing",
	},
	{
		"failed to send unregister error response",
		"authorization-token",
		"thing-id",
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{ReturnErr: errors.New("unable to unregister thing")},
		&mocks.FakePublisher{
			Err:       errors.New("unable to unregister thing"),
			SendError: errors.New("failed to send unregister response"),
		},
		&mocks.FakeConnector{},
		"failed to send unregister response",
	},
	{
		"failed to send unregister message to connector",
		"authorization-token",
		"thing-id",
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
		&mocks.FakeConnector{SendError: errors.New("failed to send unregister message")},
		"",
	},
	{
		"failed to send unregister success response",
		"authorization-token",
		"thing-id",
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{SendError: errors.New("failed to send unregister response")},
		&mocks.FakeConnector{},
		"failed to send unregister response",
	},
	{
		"allowed to unregister the thing",
		"authorization-token",
		"thing-id",
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
		&mocks.FakeConnector{},
		"",
	},
	{
		"unregister response successfully sent",
		"authorization-token",
		"thing-id",
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
		&mocks.FakeConnector{},
		"",
	},
}

func TestUnregisterThing(t *testing.T) {
	for _, tc := range unregisterAtCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.fakeLogger.
				On("Error", tc.fakeConnector.SendError).
				Maybe()
			tc.fakeThingProxy.
				On("Remove", tc.authParam, tc.idParam).
				Return(tc.fakeThingProxy.ReturnErr).
				Maybe()
			tc.fakePublisher.
				On("SendUnregisteredDevice", tc.idParam, tc.fakePublisher.Err).
				Return(tc.fakePublisher.SendError).
				Maybe()
			tc.fakeConnector.
				On("SendUnregisterDevice", tc.idParam).
				Return(tc.fakeConnector.SendError).
				Maybe()

			thingInteractor := NewThingInteractor(tc.fakeLogger, tc.fakePublisher, tc.fakeThingProxy, tc.fakeConnector)
			err := thingInteractor.Unregister(tc.authParam, tc.idParam)

			if err != nil {
				assert.EqualError(t, err, tc.expectedErrMsg)
			} else {
				assert.EqualValues(t, tc.expectedErrMsg, "")
			}

			tc.fakeLogger.AssertExpectations(t)
			tc.fakeThingProxy.AssertExpectations(t)
			tc.fakePublisher.AssertExpectations(t)
			tc.fakeConnector.AssertExpectations(t)
		})
	}
}
