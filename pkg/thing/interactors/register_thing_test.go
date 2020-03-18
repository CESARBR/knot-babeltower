package interactors

import (
	"errors"
	"testing"

	sharedEntities "github.com/CESARBR/knot-babeltower/pkg/entities"
	"github.com/CESARBR/knot-babeltower/pkg/thing/mocks"
	"github.com/stretchr/testify/assert"
)

type registerTestSuite struct {
	thingID        string
	thingName      string
	authorization  string
	fakeLogger     *mocks.FakeLogger
	fakePublisher  *mocks.FakePublisher
	fakeThingProxy *mocks.FakeThingProxy
	fakeConnector  *mocks.FakeConnector
	errExpected    error
}

func TestRegisterThing(t *testing.T) {
	testCases := map[string]registerTestSuite{
		"TestPublisherError": {
			"123",
			"test",
			"authorization token",
			&mocks.FakeLogger{},
			&mocks.FakePublisher{ReturnErr: errors.New("mock publisher error")},
			&mocks.FakeThingProxy{},
			&mocks.FakeConnector{},
			errors.New("mock publisher error"),
		},
		"TestProxyError": {
			"123",
			"test",
			"authorization token",
			&mocks.FakeLogger{},
			&mocks.FakePublisher{Token: "", SendError: errors.New("mock proxy error")},
			&mocks.FakeThingProxy{ReturnErr: errors.New("mock proxy error")},
			&mocks.FakeConnector{},
			errors.New("mock proxy error"),
		},
		"TestIDLenght": {
			"01234567890123456789",
			"test",
			"authorization token",
			&mocks.FakeLogger{},
			&mocks.FakePublisher{Token: "", SendError: ErrorIDLenght{}},
			&mocks.FakeThingProxy{},
			&mocks.FakeConnector{},
			ErrorIDLenght{},
		},
		"TestIDInvalid": {
			"not hex string",
			"test",
			"authorization token",
			&mocks.FakeLogger{},
			&mocks.FakePublisher{Token: "", SendError: ErrorIDInvalid{}},
			&mocks.FakeThingProxy{},
			&mocks.FakeConnector{},
			ErrorIDInvalid{},
		},
		"shouldRaiseConnectorSendError": {
			"123",
			"test",
			"authorization token",
			&mocks.FakeLogger{},
			&mocks.FakePublisher{},
			&mocks.FakeThingProxy{},
			&mocks.FakeConnector{SendError: sharedEntities.ErrEntityExists{}},
			sharedEntities.ErrEntityExists{},
		},
	}

	t.Logf("Number of test cases: %d", len(testCases))
	for tcName, tc := range testCases {
		t.Logf("Test case %s", tcName)
		t.Run(tcName, func(t *testing.T) {
			var err error
			var tmp *string
			if tc.fakePublisher.SendError != nil {
				tmp = new(string)
				*tmp = tc.fakePublisher.SendError.Error()
			}

			tc.fakePublisher.On("SendRegisteredDevice", tc.thingID, tc.fakePublisher.Token, tmp).
				Return(tc.fakePublisher.ReturnErr)
			tc.fakeThingProxy.On("Create", tc.thingID, tc.thingName, tc.authorization).
				Return(tc.fakePublisher.Token, tc.fakeThingProxy.ReturnErr).Maybe()
			tc.fakeConnector.On("SendRegisterDevice", tc.thingID, tc.thingName).
				Return(tc.fakeConnector.SendError).Maybe()

			thingInteractor := NewThingInteractor(tc.fakeLogger, tc.fakePublisher, tc.fakeThingProxy, tc.fakeConnector)
			err = thingInteractor.Register(tc.authorization, tc.thingID, tc.thingName)
			if err != nil && !assert.IsType(t, errors.Unwrap(err), tc.errExpected) {
				t.Errorf("Create Thing failed with unexpected error. Error: %s", err)
				return
			}

			tc.fakePublisher.AssertExpectations(t)
			tc.fakeThingProxy.AssertExpectations(t)
			tc.fakeConnector.AssertExpectations(t)
			t.Log("Create thing ok")
		})
	}
}
