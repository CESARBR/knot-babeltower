package interactors

import (
	"errors"
	"testing"

	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
	"github.com/CESARBR/knot-babeltower/pkg/thing/mocks"
	"github.com/stretchr/testify/assert"
)

type UpdateDataTestCase struct {
	name           string
	authParam      string
	idParam        string
	dataParam      []entities.Data
	fakeLogger     *mocks.FakeLogger
	fakeThingProxy *mocks.FakeThingProxy
	fakePublisher  *mocks.FakePublisher
	expectedError  error
}

var (
	errThingProxyGet = errors.New("error in thing's service")
	errClientSend    = errors.New("error sending message to client")
)

var voltageSchema = []entities.Schema{entities.Schema{
	SensorID:  0,
	ValueType: 1,
	Unit:      1,
	TypeID:    1,
	Name:      "voltage-v",
}}

var updateDataUseCases = []UpdateDataTestCase{
	{
		"authorization token not provided",
		"",
		"thing-id",
		[]entities.Data{entities.Data{}},
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
		ErrAuthNotProvided,
	},
	{
		"thing's id not provided",
		"authorization-token",
		"",
		[]entities.Data{entities.Data{}},
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
		ErrIDNotProvided,
	},
	{
		"thing's data token not provided",
		"authorization-token",
		"thing-id",
		nil,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
		ErrDataNotProvided,
	},
	{
		"failed to get thing from thing's service",
		"authorization-token",
		"thing-id",
		[]entities.Data{entities.Data{}},
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{ReturnErr: errThingProxyGet},
		&mocks.FakePublisher{},
		errThingProxyGet,
	},
	{
		"thing doesn't have a schema yet",
		"authorization-token",
		"thing-id",
		[]entities.Data{entities.Data{}},
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{Thing: &entities.Thing{
			ID:    "thing-id",
			Token: "thing-token",
			Name:  "thing",
		}},
		&mocks.FakePublisher{},
		ErrNoSchema,
	},
	{
		"data value doesn't match with thing's schema",
		"authorization-token",
		"thing-id",
		[]entities.Data{entities.Data{SensorID: 0, Value: false}},
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{Thing: &entities.Thing{
			ID:     "thing-id",
			Token:  "thing-token",
			Name:   "thing",
			Schema: voltageSchema,
		}},
		&mocks.FakePublisher{},
		ErrDataInvalid,
	},
	{
		"data sensorId doesn't match with thing's schema",
		"authorization-token",
		"thing-id",
		[]entities.Data{entities.Data{SensorID: 1, Value: 5}},
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{Thing: &entities.Thing{
			ID:     "thing-id",
			Token:  "thing-token",
			Name:   "thing",
			Schema: voltageSchema,
		}},
		&mocks.FakePublisher{},
		ErrDataInvalid,
	},
	{
		"error publishing message in client exchange",
		"authorization-token",
		"thing-id",
		[]entities.Data{entities.Data{SensorID: 0, Value: 5}},
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{Thing: &entities.Thing{
			ID:     "thing-id",
			Token:  "thing-token",
			Name:   "thing",
			Schema: voltageSchema,
		}},
		&mocks.FakePublisher{ReturnErr: errClientSend},
		errClientSend,
	},
	{
		"message successfuly send to client exchange",
		"authorization token",
		"thing-id",
		[]entities.Data{entities.Data{SensorID: 0, Value: 5}},
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{Thing: &entities.Thing{
			ID:     "thing-id",
			Token:  "thing-token",
			Name:   "thing",
			Schema: voltageSchema,
		}},
		&mocks.FakePublisher{},
		nil,
	},
}

func TestUpdateData(t *testing.T) {
	for _, tc := range updateDataUseCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.fakeThingProxy.
				On("Get", tc.authParam, tc.idParam).
				Return(tc.fakeThingProxy.Thing, tc.fakeThingProxy.ReturnErr).
				Maybe()
			tc.fakePublisher.
				On("SendUpdateData", tc.idParam, tc.dataParam).
				Return(tc.fakePublisher.ReturnErr).
				Maybe()

			thingInteractor := NewThingInteractor(tc.fakeLogger, tc.fakePublisher, tc.fakeThingProxy, nil)
			err := thingInteractor.UpdateData(tc.authParam, tc.idParam, tc.dataParam)

			assert.EqualValues(t, errors.Is(err, tc.expectedError), true)

			tc.fakeThingProxy.AssertExpectations(t)
			tc.fakePublisher.AssertExpectations(t)
		})
	}
}
