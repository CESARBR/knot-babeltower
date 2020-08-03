package interactors

import (
	"errors"
	"testing"

	"github.com/CESARBR/knot-babeltower/pkg/mocks"
	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
	"github.com/stretchr/testify/assert"
)

type UpdateConfigTestCase struct {
	name           string
	authParam      string
	idParam        string
	configParam    []entities.Config
	expectedError  error
	fakeLogger     *mocks.FakeLogger
	fakeThingProxy *mocks.FakeThingProxy
	fakePublisher  *mocks.FakePublisher
}

var configExample = []entities.Config{
	{
		SensorID:       0,
		Change:         true,
		TimeSec:        12,
		LowerThreshold: 1000,
		UpperThreshold: 2000,
	},
}

var (
	errProxyUpdateConfig = errors.New("can't update config on thing's service")
)

var updateConfigTestCases = []UpdateConfigTestCase{
	{
		"authorization token not provided",
		"",
		"c09660af89ecba61",
		[]entities.Config{{}},
		ErrAuthNotProvided,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
	},
	{
		"thing id not provided",
		"authorization-token",
		"",
		[]entities.Config{{}},
		ErrIDNotProvided,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
	},
	{
		"thing's config not provided",
		"authorization-token",
		"c09660af89ecba61",
		nil,
		ErrConfigNotProvided,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
	},
	{
		"thing's config successfully updated on the thing's proxy",
		"authorization-token",
		"c09660af89ecba61",
		configExample,
		nil,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{Thing: &entities.Thing{
			ID:     "thing-id",
			Token:  "thing-token",
			Name:   "thing",
			Schema: voltageSchema,
		}, ReturnErr: nil},
		&mocks.FakePublisher{},
	},
	{
		"failed to update thing's config on the thing's proxy",
		"authorization-token",
		"c09660af89ecba61",
		configExample,
		errProxyUpdateConfig,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{ReturnErr: errProxyUpdateConfig},
		&mocks.FakePublisher{},
	},
	{
		"failed get thing's metadata on the thing's proxy",
		"authorization-token",
		"c09660af89ecba61",
		configExample,
		errProxyUpdateConfig,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{ReturnErr: errProxyUpdateConfig},
		&mocks.FakePublisher{},
	},
	{
		"failed to updade thing's config if it doesn't have schema yet",
		"authorization-token",
		"c09660af89ecba61",
		configExample,
		ErrSchemaUndefined,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{Thing: &entities.Thing{
			ID:    "thing-id",
			Token: "thing-token",
			Name:  "thing",
		}},
		&mocks.FakePublisher{},
	},
	{
		"failed to updade thing's config if incompatible with schema",
		"authorization-token",
		"c09660af89ecba61",
		[]entities.Config{
			{
				SensorID:       1,
				Change:         true,
				TimeSec:        12,
				LowerThreshold: 1000,
				UpperThreshold: 2000,
			},
		},
		ErrConfigInvalid,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{Thing: &entities.Thing{
			ID:     "thing-id",
			Token:  "thing-token",
			Name:   "thing",
			Schema: voltageSchema,
		}},
		&mocks.FakePublisher{},
	},
}

func TestUpdateConfig(t *testing.T) {
	for _, tc := range updateConfigTestCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.fakeThingProxy.
				On("Get", tc.authParam, tc.idParam).
				Return(tc.fakeThingProxy.Thing, tc.fakeThingProxy.ReturnErr).
				Maybe()
			tc.fakeThingProxy.
				On("UpdateConfig", tc.authParam, tc.idParam, tc.configParam).
				Return(tc.fakeThingProxy.ReturnErr).
				Maybe()

			thingInteractor := NewThingInteractor(tc.fakeLogger, tc.fakePublisher, tc.fakeThingProxy)
			err := thingInteractor.UpdateConfig(tc.authParam, tc.idParam, tc.configParam)

			assert.EqualValues(t, errors.Is(err, tc.expectedError), true)
			tc.fakeThingProxy.AssertExpectations(t)
		})
	}
}
