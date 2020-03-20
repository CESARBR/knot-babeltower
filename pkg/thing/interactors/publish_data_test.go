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
	fakeConnector  *mocks.FakeConnector
	expectedError  error
}

var (
	errConnectorSend = errors.New("error sending message to connector")
)

var publishDataUseCases = []PublishDataTestCase{
	{
		"authorization token not provided",
		"",
		"thing-id",
		[]entities.Data{entities.Data{}},
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakeConnector{},
		ErrNoAuthToken,
	},
	{
		"thing's id not provided",
		"authorization-token",
		"",
		[]entities.Data{entities.Data{}},
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakeConnector{},
		ErrNoIDParam,
	},
	{
		"thing's data token not provided",
		"authorization-token",
		"thing-id",
		nil,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakeConnector{},
		ErrNoDataParam,
	},
	{
		"failed to get thing from thing's service",
		"authorization-token",
		"thing-id",
		[]entities.Data{entities.Data{}},
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{ReturnErr: errThingProxyGet},
		&mocks.FakeConnector{},
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
		&mocks.FakeConnector{},
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
		&mocks.FakeConnector{},
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
		&mocks.FakeConnector{},
		ErrDataInvalid,
	},
	{
		"error publishing message in connector exchange",
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
		&mocks.FakeConnector{SendError: errConnectorSend},
		errConnectorSend,
	},
	{
		"message successfully sent to connector exchange",
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
		&mocks.FakeConnector{},
		nil,
	},
}

func TestPublishData(t *testing.T) {
	for _, tc := range publishDataUseCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.fakeThingProxy.
				On("Get", tc.authParam, tc.idParam).
				Return(tc.fakeThingProxy.Thing, tc.fakeThingProxy.ReturnErr).
				Maybe()
			tc.fakeConnector.
				On("SendPublishData", tc.idParam, tc.dataParam).
				Return(tc.fakeConnector.SendError).
				Maybe()

			thingInteractor := NewThingInteractor(tc.fakeLogger, nil, tc.fakeThingProxy, tc.fakeConnector)
			err := thingInteractor.PublishData(tc.authParam, tc.idParam, tc.dataParam)

			assert.EqualValues(t, errors.Is(err, tc.expectedError), true)

			tc.fakeThingProxy.AssertExpectations(t)
			tc.fakeConnector.AssertExpectations(t)
		})
	}
}
