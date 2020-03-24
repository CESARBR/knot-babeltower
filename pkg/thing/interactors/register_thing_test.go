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
	errExpected    error
	fakeLogger     *mocks.FakeLogger
	fakePublisher  *mocks.FakePublisher
	fakeThingProxy *mocks.FakeThingProxy
	fakeConnector  *mocks.FakeConnector
}

func TestRegisterThing(t *testing.T) {
	testCases := map[string]registerTestSuite{
		"TestPublisherError": {
			"123",
			"test",
			"authorization token",
			errors.New("mock publisher error"),
			&mocks.FakeLogger{},
			&mocks.FakePublisher{ReturnErr: errors.New("mock publisher error")},
			&mocks.FakeThingProxy{},
			&mocks.FakeConnector{},
		},
		"TestProxyError": {
			"123",
			"test",
			"authorization token",
			errors.New("mock proxy error"),
			&mocks.FakeLogger{},
			&mocks.FakePublisher{Token: "", SendError: errors.New("mock proxy error")},
			&mocks.FakeThingProxy{ReturnErr: errors.New("mock proxy error")},
			&mocks.FakeConnector{},
		},
		"TestIDLenght": {
			"01234567890123456789",
			"test",
			"authorization token",
			ErrIDLength,
			&mocks.FakeLogger{},
			&mocks.FakePublisher{Token: "", SendError: ErrIDLength},
			&mocks.FakeThingProxy{},
			&mocks.FakeConnector{},
		},
		"TestIDInvalid": {
			"not hex string",
			"test",
			"authorization token",
			ErrIDNotHex,
			&mocks.FakeLogger{},
			&mocks.FakePublisher{Token: "", SendError: ErrIDNotHex},
			&mocks.FakeThingProxy{},
			&mocks.FakeConnector{},
		},
		"shouldRaiseConnectorSendError": {
			"123",
			"test",
			"authorization token",
			sharedEntities.ErrEntityExists{},
			&mocks.FakeLogger{},
			&mocks.FakePublisher{},
			&mocks.FakeThingProxy{},
			&mocks.FakeConnector{SendError: sharedEntities.ErrEntityExists{}},
		},
	}

	t.Logf("Number of test cases: %d", len(testCases))
	for tcName, tc := range testCases {
		t.Logf("Test case %s", tcName)
		t.Run(tcName, func(t *testing.T) {
			tc.fakePublisher.On("SendRegisteredDevice", tc.thingID, tc.fakePublisher.Token, tc.fakePublisher.SendError).
				Return(tc.fakePublisher.ReturnErr)
			tc.fakeThingProxy.On("Create", tc.thingID, tc.thingName, tc.authorization).
				Return(tc.fakePublisher.Token, tc.fakeThingProxy.ReturnErr).Maybe()
			tc.fakeConnector.On("SendRegisterDevice", tc.thingID, tc.thingName).
				Return(tc.fakeConnector.SendError).Maybe()

			thingInteractor := NewThingInteractor(tc.fakeLogger, tc.fakePublisher, tc.fakeThingProxy, tc.fakeConnector)
			err := thingInteractor.Register(tc.authorization, tc.thingID, tc.thingName)
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
