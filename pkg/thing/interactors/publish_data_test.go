package interactors

import (
	"errors"
	"testing"

	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
	"github.com/CESARBR/knot-babeltower/pkg/thing/mocks"
	"github.com/stretchr/testify/assert"
)

type PublishDataTestCase struct {
	name           string
	authParam      string
	idParam        string
	dataParam      []entities.Data
	fakeLogger     *mocks.FakeLogger
	fakeThingProxy *mocks.FakeThingProxy
	expectedError  error
}

var publishDataUseCases = []PublishDataTestCase{
	{
		"authorization token not provided",
		"",
		"thing-id",
		[]entities.Data{entities.Data{}},
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		ErrAuthNotProvided,
	},
	{
		"thing's id not provided",
		"authorization-token",
		"",
		[]entities.Data{entities.Data{}},
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		ErrIDNotProvided,
	},
	{
		"thing's data token not provided",
		"authorization-token",
		"thing-id",
		nil,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		ErrDataNotProvided,
	},
	{
		"failed to get thing from thing's service",
		"authorization-token",
		"thing-id",
		[]entities.Data{entities.Data{}},
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{ReturnErr: errThingProxyGet},
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
		ErrSchemaUndefined,
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
		ErrDataInvalid,
	},
	{
		"data sensorId doesn't match with thing's schema",
		"authorization-token",
		"thing-id",
		[]entities.Data{entities.Data{SensorID: 1, Value: float64(5)}},
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{Thing: &entities.Thing{
			ID:     "thing-id",
			Token:  "thing-token",
			Name:   "thing",
			Schema: voltageSchema,
		}},
		ErrDataInvalid,
	},
}

func TestPublishData(t *testing.T) {
	for _, tc := range publishDataUseCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.fakeThingProxy.
				On("Get", tc.authParam, tc.idParam).
				Return(tc.fakeThingProxy.Thing, tc.fakeThingProxy.ReturnErr).
				Maybe()

			thingInteractor := NewThingInteractor(tc.fakeLogger, nil, tc.fakeThingProxy)
			err := thingInteractor.PublishData(tc.authParam, tc.idParam, tc.dataParam)

			assert.EqualValues(t, errors.Is(err, tc.expectedError), true)

			tc.fakeThingProxy.AssertExpectations(t)
		})
	}
}
