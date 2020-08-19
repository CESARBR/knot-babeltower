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
		LowerThreshold: 25.4,
		UpperThreshold: 87.2,
	},
	{
		SensorID:       1,
		Change:         false,
		TimeSec:        60,
		LowerThreshold: 25.4,
	},
}

var schemaExample = []entities.Schema{
	{
		SensorID:  3,
		ValueType: 1,
		Unit:      0,
		TypeID:    65521,
		Name:      "thing-with-float-data",
	},
	{
		SensorID:  0,
		ValueType: 2,
		Unit:      0,
		TypeID:    65521,
		Name:      "thing-with-float-data",
	},
	{
		SensorID:  1,
		ValueType: 2,
		Unit:      0,
		TypeID:    65521,
		Name:      "second-thing-with-float-data",
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
			Schema: schemaExample,
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
		"failed to updade thing's config if has no schema associated with the same sensorId",
		"authorization-token",
		"c09660af89ecba61",
		[]entities.Config{{
			SensorID:       10, // different sensorId
			Change:         true,
			TimeSec:        12,
			LowerThreshold: 25.4,
			UpperThreshold: 87.2,
		}},
		ErrConfigInvalid,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{Thing: &entities.Thing{
			ID:     "thing-id",
			Token:  "thing-token",
			Name:   "thing",
			Schema: schemaExample,
		}},
		&mocks.FakePublisher{},
	},
	{
		"failed to updade thing's config if threshold value is incompatible with schema's valueType",
		"authorization-token",
		"c09660af89ecba61",
		[]entities.Config{{
			SensorID:       0,
			Change:         true,
			TimeSec:        12,
			LowerThreshold: 400, // incompatible with valueType 2: floats.
			UpperThreshold: 87.2,
		}},
		ErrDataInvalid,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{Thing: &entities.Thing{
			ID:     "thing-id",
			Token:  "thing-token",
			Name:   "thing",
			Schema: schemaExample,
		}},
		&mocks.FakePublisher{},
	},
	{
		"failed to updade thing's config if it has not changed",
		"authorization-token",
		"c09660af89ecba61",
		configExample,
		ErrConfigEqual,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{Thing: &entities.Thing{
			ID:     "thing-id",
			Token:  "thing-token",
			Name:   "thing",
			Schema: schemaExample,
			Config: configExample,
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
