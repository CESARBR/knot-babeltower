package interactors

import (
	"errors"
	"testing"

	"github.com/CESARBR/knot-babeltower/pkg/mocks"
	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
	"github.com/stretchr/testify/assert"
)

type PublishDataTestCase struct {
	name             string
	authParam        string
	idParam          string
	dataParam        []entities.Data
	fakeLogger       *mocks.FakeLogger
	fakeThingProxy   *mocks.FakeThingProxy
	fakePublisher    *mocks.FakePublisher
	fakeSessionStore *mocks.FakeSessionStore
	expectedError    error
}

var (
	errPublishData        = errors.New("error publishing data in broadcast mode")
	errPublishSessionData = errors.New("error publishing data to user sessions")
	tokenWithValidEmail   = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTI0NTE3OTEsImp0aSI6IjgwYjBlNzk5LTAzNjItNGE0NC1iMTA3LWMwOTNjYzRiMjM2MSIsImlhdCI6MTYyMDkxNTc5MSwiaXNzIjoiamFzbkBjZXNhci5vcmcuYnIiLCJ0eXBlIjoyfQ._TBgxNAf18f5_FdH1oCuYe1v3NPOyL68l0-nzx4XAI8"
	emailExample          = "jasn@cesar.org.br"
	errGetSession         = errors.New("error getting user session")
)

var publishDataUseCases = []PublishDataTestCase{
	{
		"authorization token not provided",
		"",
		"thing-id",
		[]entities.Data{{}},
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
		&mocks.FakeSessionStore{},
		ErrAuthNotProvided,
	},
	{
		"thing's id not provided",
		"authorization-token",
		"",
		[]entities.Data{{}},
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
		&mocks.FakeSessionStore{},
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
		&mocks.FakeSessionStore{},
		ErrDataNotProvided,
	},
	{
		"failed to get thing from thing's service",
		"authorization-token",
		"thing-id",
		[]entities.Data{{}},
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{ReturnErr: errThingProxyGet},
		&mocks.FakePublisher{},
		&mocks.FakeSessionStore{},
		errThingProxyGet,
	},
	{
		"thing doesn't have a config yet",
		"authorization-token",
		"thing-id",
		[]entities.Data{{}},
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{Thing: &entities.Thing{
			ID:    "thing-id",
			Token: "thing-token",
			Name:  "thing",
		}},
		&mocks.FakePublisher{},
		&mocks.FakeSessionStore{},
		ErrConfigUndefined,
	},
	{
		"data value doesn't match with thing's schema",
		"authorization-token",
		"thing-id",
		[]entities.Data{{SensorID: 0, Value: false}},
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{Thing: &entities.Thing{
			ID:     "thing-id",
			Token:  "thing-token",
			Name:   "thing",
			Config: configWithVoltageSchema,
		}},
		&mocks.FakePublisher{},
		&mocks.FakeSessionStore{},
		ErrDataInvalid,
	},
	{
		"data sensorId doesn't match with thing's schema",
		"authorization-token",
		"thing-id",
		[]entities.Data{{SensorID: 1, Value: float64(5)}},
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{Thing: &entities.Thing{
			ID:     "thing-id",
			Token:  "thing-token",
			Name:   "thing",
			Config: configWithVoltageSchema,
		}},
		&mocks.FakePublisher{},
		&mocks.FakeSessionStore{},
		ErrDataInvalid,
	},
	{
		"failed to publish broadcast data",
		"authorization-token",
		"thing-id",
		[]entities.Data{{SensorID: 0, Value: float64(5)}},
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{Thing: &entities.Thing{
			ID:     "thing-id",
			Token:  "thing-token",
			Name:   "thing",
			Config: configWithVoltageSchema,
		}},
		&mocks.FakePublisher{PublishErr: errPublishData},
		&mocks.FakeSessionStore{},
		errPublishData,
	},
	{
		"failed to get user session",
		tokenWithValidEmail,
		"thing-id",
		[]entities.Data{{SensorID: 0, Value: float64(5)}},
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{Thing: &entities.Thing{
			ID:     "thing-id",
			Token:  "thing-token",
			Name:   "thing",
			Config: configWithVoltageSchema,
		}},
		&mocks.FakePublisher{PublishSessionErr: errPublishSessionData},
		&mocks.FakeSessionStore{GetReturnErr: errGetSession},
		errGetSession,
	},
	{
		"failed to publish session data",
		tokenWithValidEmail,
		"thing-id",
		[]entities.Data{{SensorID: 0, Value: float64(5)}},
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{Thing: &entities.Thing{
			ID:     "thing-id",
			Token:  "thing-token",
			Name:   "thing",
			Config: configWithVoltageSchema,
		}},
		&mocks.FakePublisher{PublishSessionErr: errPublishSessionData},
		&mocks.FakeSessionStore{Session: "session-id"},
		errPublishSessionData,
	},
	{
		"data successfully published",
		tokenWithValidEmail,
		"thing-id",
		[]entities.Data{{SensorID: 0, Value: float64(5)}},
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{Thing: &entities.Thing{
			ID:     "thing-id",
			Token:  "thing-token",
			Name:   "thing",
			Config: configWithVoltageSchema,
		}},
		&mocks.FakePublisher{},
		&mocks.FakeSessionStore{Session: "session-id"},
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
			tc.fakePublisher.
				On("PublishBroadcastData", tc.idParam, tc.authParam, tc.dataParam).
				Return(tc.fakePublisher.PublishErr).
				Maybe()
			tc.fakePublisher.
				On("PublishSessionData", tc.idParam, tc.authParam, tc.fakeSessionStore.Session, tc.dataParam).
				Return(tc.fakePublisher.PublishSessionErr).
				Maybe()
			tc.fakeSessionStore.
				On("Get", emailExample).
				Return(tc.fakeSessionStore.Session, tc.fakeSessionStore.GetReturnErr).
				Maybe()

			thingInteractor := NewThingInteractor(tc.fakeLogger, tc.fakePublisher, tc.fakeThingProxy, tc.fakeSessionStore)
			err := thingInteractor.PublishData(tc.authParam, tc.idParam, tc.dataParam)
			assert.EqualValues(t, errors.Is(err, tc.expectedError), true)

			tc.fakeThingProxy.AssertExpectations(t)
			tc.fakePublisher.AssertExpectations(t)
			tc.fakeSessionStore.AssertExpectations(t)
		})
	}
}
