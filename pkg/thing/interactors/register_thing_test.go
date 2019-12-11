package interactors

import (
	"errors"
	"testing"

	sharedEntities "github.com/CESARBR/knot-babeltower/pkg/entities"
	"github.com/CESARBR/knot-babeltower/pkg/network"
	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type registerTestSuite struct {
	thingID       string
	thingName     string
	authorization string
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
	sendError error
	recvError error
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

func (fp *FakeMsgPublisher) SendUpdatedSchema(thingID string) error {
	ret := fp.Called(thingID)

	return ret.Error(0)
}

func (fp *FakeProxy) Create(id, name, authorization string) (string, error) {
	ret := fp.Called(id, name, authorization)

	return ret.String(0), ret.Error(1)
}

func (fp *FakeProxy) UpdateSchema(authorization, id string, schemaList []entities.Schema) error {
	ret := fp.Called(authorization, id, schemaList)

	return ret.Error(0)
}

func (fc *FakeConnector) SendRegisterDevice(id, name string) (err error) {
	ret := fc.Called(id, name)

	return ret.Error(0)
}

func (fc *FakeConnector) RecvRegisterDevice() (bytes []byte, err error) {
	ret := fc.Called()

	return bytes, ret.Error(1)
}

func TestRegisterThing(t *testing.T) {
	testCases := map[string]registerTestSuite{
		"TestPublisherError": {
			"123",
			"test",
			"authorization token",
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{returnErr: errors.New("mock publisher error")},
			&FakeProxy{},
			&FakeConnector{},
			errors.New("mock publisher error"),
		},
		"TestProxyError": {
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
			"01234567890123456789",
			"test",
			"authorization token",
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{token: "", sendError: ErrorIDLenght{}},
			&FakeProxy{},
			&FakeConnector{},
			ErrorIDLenght{},
		},
		"TestIDInvalid": {
			"not hex string",
			"test",
			"authorization token",
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{token: "", sendError: ErrorIDInvalid{}},
			&FakeProxy{},
			&FakeConnector{},
			ErrorIDInvalid{},
		},
		"shouldRaiseConnectorSendError": {
			"123",
			"test",
			"authorization token",
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{},
			&FakeProxy{},
			&FakeConnector{sendError: sharedEntities.ErrEntityExists{}},
			sharedEntities.ErrEntityExists{},
		},
		"shouldRaiseConnectorRecvError": {
			"123",
			"test",
			"authorization token",
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{},
			&FakeProxy{},
			&FakeConnector{recvError: sharedEntities.ErrEntityExists{}},
			sharedEntities.ErrEntityExists{},
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
				Return(tc.fakeConnector.sendError).Maybe()
			tc.fakeConnector.On("RecvRegisterDevice").
				Return([]byte{}, tc.fakeConnector.recvError).Maybe()

			thingInteractor := NewThingInteractor(tc.fakeLogger, tc.fakePublisher, tc.fakeProxy, tc.fakeConnector)
			err = thingInteractor.Register(tc.authorization, tc.thingID, tc.thingName)
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
