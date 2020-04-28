package interactors

import (
	"errors"
	"testing"

	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
	"github.com/CESARBR/knot-babeltower/pkg/thing/mocks"
	"github.com/stretchr/testify/assert"
)

type registerTestCase struct {
	name           string
	authParam      string
	idParam        string
	nameParam      string
	errExpected    error
	thingExpected  *entities.Thing
	fakeLogger     *mocks.FakeLogger
	fakeThingProxy *mocks.FakeThingProxy
	fakePublisher  *mocks.FakePublisher
	fakeConnector  *mocks.FakeConnector
}

var (
	errRegisterResponse = errors.New("error sending response to client")
	errThingCreation    = errors.New("error in thing's service")
	thing               = &entities.Thing{
		ID:    "fc3fcf912d0c290a",
		Token: "authorization-token",
		Name:  "thing",
	}
)

var registerThingUseCases = []registerTestCase{
	{
		"thing's ID has wrong lenght",
		"authorization-token",
		"01234567890123456789",
		"knot-thing",
		ErrIDLength,
		nil,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{SendError: ErrIDLength},
		&mocks.FakeConnector{},
	},
	{
		"thing's ID isn't hexadecimal",
		"authorization-token",
		"not hex string",
		"test",
		ErrIDNotHex,
		nil,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{SendError: ErrIDNotHex},
		&mocks.FakeConnector{},
	},
	{
		"thing already registered on thing's service",
		"authorization-token",
		"fc3fcf912d0c290a",
		"test",
		entities.ErrThingExists,
		thing,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{SendError: entities.ErrThingExists},
		&mocks.FakeConnector{SendError: entities.ErrThingExists},
	},
	{
		"failed to create a thing on the thing's service",
		"authorization-token",
		"fc3fcf912d0c290a",
		"knot-thing",
		errThingCreation,
		nil,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{ReturnErr: entities.ErrThingNotFound, CreateErr: errThingCreation},
		&mocks.FakePublisher{Token: "", SendError: errThingCreation},
		&mocks.FakeConnector{},
	},
	{
		"thing successfully created on thing's service",
		"authorization-token",
		"fc3fcf912d0c290a",
		"knot-thing",
		nil,
		nil,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{ReturnErr: entities.ErrThingNotFound},
		&mocks.FakePublisher{},
		&mocks.FakeConnector{},
	},
	{
		"failed to send register response",
		"authorization-token",
		"fc3fcf912d0c290a",
		"knot-thing",
		errRegisterResponse,
		nil,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{ReturnErr: entities.ErrThingNotFound},
		&mocks.FakePublisher{ReturnErr: errRegisterResponse},
		&mocks.FakeConnector{},
	},
	{
		"register response successfully sent",
		"authorization-token",
		"fc3fcf912d0c290a",
		"knot-thing",
		nil,
		nil,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{ReturnErr: entities.ErrThingNotFound},
		&mocks.FakePublisher{},
		&mocks.FakeConnector{},
	},
}

func TestRegisterThing(t *testing.T) {
	for _, tc := range registerThingUseCases {
		t.Logf("Running Test Casee: %s", tc.name)
		t.Run(tc.name, func(t *testing.T) {
			tc.fakeThingProxy.On("Get", tc.authParam, tc.idParam).
				Return(tc.thingExpected, tc.fakeThingProxy.ReturnErr).Maybe()
			tc.fakePublisher.On("SendRegisteredDevice", tc.idParam, tc.nameParam, tc.fakePublisher.Token, tc.fakePublisher.SendError).
				Return(tc.fakePublisher.ReturnErr).Maybe()
			tc.fakeThingProxy.On("Create", tc.idParam, tc.nameParam, tc.authParam).
				Return(tc.fakePublisher.Token, tc.fakeThingProxy.CreateErr).Maybe()
			tc.fakeConnector.On("SendRegisterDevice", tc.idParam, tc.nameParam).
				Return(tc.fakeConnector.SendError).Maybe()

			thingInteractor := NewThingInteractor(tc.fakeLogger, tc.fakePublisher, tc.fakeThingProxy, tc.fakeConnector)
			err := thingInteractor.Register(tc.authParam, tc.idParam, tc.nameParam)
			if err != nil && !assert.IsType(t, errors.Unwrap(err), tc.errExpected) {
				t.Errorf("create thing failed with unexpected error. Error: %s", err)
				return
			}

			tc.fakePublisher.AssertExpectations(t)
			tc.fakeThingProxy.AssertExpectations(t)
			tc.fakeConnector.AssertExpectations(t)
			t.Log("create thing ok")
		})
	}
}
