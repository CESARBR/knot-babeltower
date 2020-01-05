package interactors

import (
	"errors"
	"testing"

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
	errExpected   error
}

type FakeRegisterThingLogger struct {
}

type FakeMsgPublisher struct {
	mock.Mock
	sendError error
	returnErr error
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

func TestRegisterThing(t *testing.T) {
	testCases := map[string]registerTestSuite{
		"TestPublisherError": {
			false,
			"123",
			"test",
			"authorization token",
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{returnErr: errors.New("mock publisher error")},
			errors.New("mock publisher error"),
		},
		"TestIDLenght": {
			false,
			"01234567890123456789",
			"test",
			"authorization token",
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{sendError: ErrorIDLenght{}},
			ErrorIDLenght{},
		},
		"TestNameEmpty": {
			false,
			"123",
			"",
			"authorization token",
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{sendError: ErrorNameNotFound{}},
			ErrorNameNotFound{},
		},
		"TestIDInvalid": {
			false,
			"not hex string",
			"test",
			"authorization token",
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{sendError: ErrorIDInvalid{}},
			ErrorIDInvalid{},
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

			msg := network.RegisterResponseMsg{ID: tc.thingID, Token: "", Error: tmp}
			tc.fakePublisher.On("SendRegisterDevice", msg).
				Return(tc.fakePublisher.returnErr)
			createThingInteractor := NewRegisterThing(tc.fakeLogger, tc.fakePublisher)
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
			t.Log("Create thing ok")
		})
	}
}
