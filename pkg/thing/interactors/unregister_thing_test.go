package interactors

import (
	"errors"
	"testing"

	"github.com/CESARBR/knot-babeltower/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

type UnregisterThingTestCase struct {
	name           string
	authParam      string
	idParam        string
	fakeLogger     *mocks.FakeLogger
	fakeThingProxy *mocks.FakeThingProxy
	fakePublisher  *mocks.FakePublisher
	expectedErrMsg string
}

var unregisterAtCases = []UnregisterThingTestCase{
	{
		"authorization token not provided",
		"",
		"thing-id",
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
		ErrAuthNotProvided.Error(),
	},
	{
		"thing's id not provided",
		"authorization-token",
		"",
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
		ErrIDNotProvided.Error(),
	},
	{
		"thing's id not found",
		"authorization-token",
		"thing-id",
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{ReturnErr: errors.New("thing's id not found")},
		&mocks.FakePublisher{Err: errors.New("thing's id not found")},
		"thing's id not found",
	},
	{
		"unable to unregister thing",
		"authorization-token",
		"thing-id",
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{ReturnErr: errors.New("unable to unregister thing")},
		&mocks.FakePublisher{Err: errors.New("unable to unregister thing")},
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
		"failed to send unregister response",
	},
	{
		"failed to send unregister success response",
		"authorization-token",
		"thing-id",
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{SendError: errors.New("failed to send unregister response")},
		"failed to send unregister response",
	},
	{
		"allowed to unregister the thing",
		"authorization-token",
		"thing-id",
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
		"",
	},
	{
		"unregister response successfully sent",
		"authorization-token",
		"thing-id",
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
		"",
	},
}

func TestUnregisterThing(t *testing.T) {
	for _, tc := range unregisterAtCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.fakeThingProxy.
				On("Remove", tc.authParam, tc.idParam).
				Return(tc.fakeThingProxy.ReturnErr).
				Maybe()
			tc.fakePublisher.
				On("PublishUnregisteredDevice", tc.idParam, tc.fakePublisher.Err).
				Return(tc.fakePublisher.SendError).
				Maybe()

			thingInteractor := NewThingInteractor(tc.fakeLogger, tc.fakePublisher, tc.fakeThingProxy)
			err := thingInteractor.Unregister(tc.authParam, tc.idParam)

			if err != nil {
				assert.EqualError(t, err, tc.expectedErrMsg)
			} else {
				assert.EqualValues(t, tc.expectedErrMsg, "")
			}

			tc.fakeLogger.AssertExpectations(t)
			tc.fakeThingProxy.AssertExpectations(t)
			tc.fakePublisher.AssertExpectations(t)
		})
	}
}
