package interactors

import (
	"errors"
	"testing"

	"github.com/CESARBR/knot-babeltower/pkg/entities"
	"github.com/CESARBR/knot-babeltower/pkg/network"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type registerTestSuite struct {
	testArguments bool
	thingID       string
	thingName     interface{}
	authorization interface{}
	fakeLogger    *FakeRegisterThingLogger
	fakePublisher *FakeMsgPublisher
	fakeProxy     *FakeProxy
	fakeConnector *FakeConnector
	errExpected   error
}

type FakeRegisterThingLogger struct {
}

type FakeMsgPublisher struct {
	mock.Mock
	sendError error
	returnErr error
	token     string
}

type FakeProxy struct {
	mock.Mock
	returnError error
}

type FakeConnector struct {
	mock.Mock
	returnError error
}

func (fl *FakeRegisterThingLogger) Info(...interface{}) {}

func (fl *FakeRegisterThingLogger) Infof(string, ...interface{}) {}

func (fl *FakeRegisterThingLogger) Debug(...interface{}) {}

func (fl *FakeRegisterThingLogger) Warn(...interface{}) {}

func (fl *FakeRegisterThingLogger) Error(...interface{}) {}

func (fl *FakeRegisterThingLogger) Errorf(string, ...interface{}) {}

func (fp *FakeMsgPublisher) SendRegisterDevice(msg network.RegisterResponseMsg) error {
	ret := fp.Called(msg)

	return ret.Error(0)
}

func (fp *FakeProxy) Create(id, name, authorization string) (string, error) {
	ret := fp.Called(id, name, authorization)

	return ret.String(0), ret.Error(1)
}

func (fc *FakeConnector) SendRegisterDevice(id, name string) (err error) {
	ret := fc.Called(id, name)

	return ret.Error(0)
}

func TestRegisterThing(t *testing.T) {
	testCases := map[string]registerTestSuite{
		"TestPublisherError": {
			false,
			"123",
			"test",
			"authorization token",
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{returnErr: errors.New("mock publisher error")},
			&FakeProxy{},
			&FakeConnector{returnError: nil},
			errors.New("mock publisher error"),
		},
		"TestProxyError": {
			false,
			"123",
			"test",
			"authorization token",
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{token: "", sendError: errors.New("mock proxy error")},
			&FakeProxy{returnError: errors.New("mock proxy error")},
			&FakeConnector{},
			errors.New("mock proxy error"),
		},
		"TestIDLenght": {
			false,
			"01234567890123456789",
			"test",
			"authorization token",
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{token: "", sendError: ErrorIDLenght{}},
			&FakeProxy{},
			&FakeConnector{},
			ErrorIDLenght{},
		},
		"TestNameEmpty": {
			false,
			"123",
			"",
			"authorization token",
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{token: "", sendError: ErrorNameNotFound{}},
			&FakeProxy{},
			&FakeConnector{},
			ErrorNameNotFound{},
		},
		"TestIDInvalid": {
			false,
			"not hex string",
			"test",
			"authorization token",
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{token: "", sendError: ErrorIDInvalid{}},
			&FakeProxy{},
			&FakeConnector{},
			ErrorIDInvalid{},
		},
		"TestMissingArgument": {
			true,
			"123",
			"",
			"",
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{token: "", sendError: ErrorMissingArgument{}},
			&FakeProxy{},
			&FakeConnector{},
			ErrorMissingArgument{},
		},
		"TestInvalidTypeName": {
			false,
			"123",
			123,
			"",
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{token: "", sendError: ErrorInvalidTypeArgument{"Name is not string"}},
			&FakeProxy{},
			&FakeConnector{},
			ErrorInvalidTypeArgument{"Name is not string"},
		},
		"TestInvalidTypeToken": {
			false,
			"123",
			"test",
			123,
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{token: "", sendError: ErrorInvalidTypeArgument{"Authorization token is not string"}},
			&FakeProxy{},
			&FakeConnector{},
			ErrorInvalidTypeArgument{"Authorization token is not string"},
		},
		"TestTokenUnauthorized": {
			false,
			"123",
			"test",
			"",
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{token: "", sendError: ErrorUnauthorized{}},
			&FakeProxy{},
			&FakeConnector{},
			ErrorUnauthorized{},
		},
		"shouldRaiseConnectorError": {
			false,
			"123",
			"test",
			"authorization token",
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{},
			&FakeProxy{},
			&FakeConnector{returnError: entities.ErrEntityExists{}},
			entities.ErrEntityExists{},
		},
	}

	t.Logf("Number of test cases: %d", len(testCases))
	for tcName, tc := range testCases {
		t.Logf("Test case %s", tcName)
		t.Run(tcName, func(t *testing.T) {
			var err error
			var tmp *string
			if tc.fakePublisher.sendError != nil {
				tmp = new(string)
				*tmp = tc.fakePublisher.sendError.Error()
			}

			msg := network.RegisterResponseMsg{ID: tc.thingID, Token: tc.fakePublisher.token, Error: tmp}
			tc.fakePublisher.On("SendRegisterDevice", msg).
				Return(tc.fakePublisher.returnErr)
			tc.fakeProxy.On("Create", tc.thingID, tc.thingName, tc.authorization).
				Return(tc.fakePublisher.token, tc.fakeProxy.returnError).Maybe()
			tc.fakeConnector.On("SendRegisterDevice", tc.thingID, tc.thingName).
				Return(tc.fakeConnector.returnError).Maybe()

			createThingInteractor := NewRegisterThing(tc.fakeLogger, tc.fakePublisher, tc.fakeProxy, tc.fakeConnector)
			if tc.testArguments {
				err = createThingInteractor.Execute(tc.thingID)
			} else {
				err = createThingInteractor.Execute(tc.thingID, tc.thingName, tc.authorization)
			}

			if err != nil && !assert.IsType(t, err, tc.errExpected) {
				t.Errorf("Create Thing failed with unexpected error. Error: %s", err)
				return
			}

			tc.fakePublisher.AssertExpectations(t)
			tc.fakeProxy.AssertExpectations(t)
			tc.fakeConnector.AssertExpectations(t)
			t.Log("Create thing ok")
		})
	}
}
