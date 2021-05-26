package interactors

import (
	"errors"
	"testing"

	"github.com/CESARBR/knot-babeltower/pkg/mocks"
	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
	"github.com/stretchr/testify/assert"
)

type UpdateConfigTestCase struct {
	name            string
	authParam       string
	idParam         string
	configParam     []entities.Config
	expectedError   error
	expectedChanged bool
	fakeLogger      *mocks.FakeLogger
	fakeThingProxy  *mocks.FakeThingProxy
	fakePublisher   *mocks.FakePublisher
}

var configExample = []entities.Config{
	{
		SensorID: 0,
		Schema: entities.Schema{
			ValueType: 2,
			Unit:      0,
			TypeID:    65521,
			Name:      "thing-with-float-data",
		},
		Event: entities.Event{
			Change:         true,
			TimeSec:        12,
			LowerThreshold: 25.4,
			UpperThreshold: 87.2,
		},
	},
	{
		SensorID: 1,
		Schema: entities.Schema{
			ValueType: 2,
			Unit:      0,
			TypeID:    65521,
			Name:      "second-thing-with-float-data",
		},
		Event: entities.Event{
			Change:         false,
			TimeSec:        60,
			LowerThreshold: 25.4,
		},
	},
	{
		SensorID: 3,
		Schema: entities.Schema{
			ValueType: 1,
			Unit:      0,
			TypeID:    65521,
			Name:      "thing-with-float-data",
		},
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
		false,
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
		false,
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
		false,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
	},
	{
		"thing's config successfully updated on the thing's proxy",
		"authorization-token",
		"c09660af89ecba61",
		[]entities.Config{{
			SensorID: 1,
			Schema: entities.Schema{
				ValueType: 2,
				Unit:      0,
				TypeID:    65521,
				Name:      "second-thing-with-float-data",
			},
			Event: entities.Event{
				Change:         false,
				TimeSec:        20,
				LowerThreshold: 25.4,
			},
		}},
		nil,
		true,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{Thing: &entities.Thing{
			ID:     "thing-id",
			Token:  "thing-token",
			Name:   "thing",
			Config: configExample,
		}, ReturnErr: nil},
		&mocks.FakePublisher{},
	},
	{
		"failed to update thing's config on the thing's proxy",
		"authorization-token",
		"c09660af89ecba61",
		configExample,
		errProxyUpdateConfig,
		false,
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
		false,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{ReturnErr: errProxyUpdateConfig},
		&mocks.FakePublisher{},
	},
	{
		"failed to updade thing's config if threshold value is incompatible with schema's valueType",
		"authorization-token",
		"c09660af89ecba61",
		[]entities.Config{{
			SensorID: 0,
			Schema: entities.Schema{
				ValueType: 2,
				Unit:      0,
				TypeID:    65521,
				Name:      "thing-with-float-data",
			},
			Event: entities.Event{
				Change:         true,
				TimeSec:        12,
				LowerThreshold: 400,
				UpperThreshold: 87.2,
			},
		}},
		ErrDataInvalid,
		false,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{Thing: &entities.Thing{
			ID:     "thing-id",
			Token:  "thing-token",
			Name:   "thing",
			Config: configExample,
		}},
		&mocks.FakePublisher{},
	},
	{
		"successfully updade thing's config if valuetype changed",
		"authorization-token",
		"c09660af89ecba61",
		[]entities.Config{{
			SensorID: 0,
			Schema: entities.Schema{
				ValueType: 3,
				Unit:      0,
				TypeID:    65521,
				Name:      "thing-with-float-data",
			},
			Event: entities.Event{
				Change:         true,
				TimeSec:        12,
				LowerThreshold: 25.4,
				UpperThreshold: 87.2,
			},
		}},
		nil,
		true,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{Thing: &entities.Thing{
			ID:    "thing-id",
			Token: "thing-token",
			Name:  "thing",
			Config: []entities.Config{
				{
					SensorID: 0,
					Schema: entities.Schema{
						ValueType: 2,
						Unit:      0,
						TypeID:    65521,
						Name:      "thing-with-float-data",
					},
					Event: entities.Event{
						Change:         true,
						TimeSec:        12,
						LowerThreshold: 25.4,
						UpperThreshold: 87.2,
					},
				},
			},
		}},
		&mocks.FakePublisher{},
	},
	{
		"failed to updade thing's sensor **event** if it hasn't a schema",
		"authorization-token",
		"c09660af89ecba61",
		[]entities.Config{{
			SensorID: 0,
			Event: entities.Event{
				Change:         true,
				TimeSec:        12,
				LowerThreshold: 25.4,
				UpperThreshold: 87.2,
			},
		}},
		ErrSchemaNotProvided,
		false,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{Thing: &entities.Thing{
			ID:     "thing-id",
			Token:  "thing-token",
			Name:   "thing",
			Config: []entities.Config{},
		}},
		&mocks.FakePublisher{},
	},
	{
		"invalid schema name",
		"authorization-token",
		"89cf40c23012ce1c",
		[]entities.Config{{
			SensorID: 0,
			Schema: entities.Schema{
				ValueType: 3,
				Unit:      0,
				TypeID:    65521,
				Name:      "SchemaNameGreaterThan23Characters",
			},
			Event: entities.Event{
				Change:         true,
				TimeSec:        12,
				LowerThreshold: 25.4,
				UpperThreshold: 87.2,
			},
		}},
		ErrSchemaInvalid,
		false,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{Thing: &entities.Thing{
			ID:     "thing-id",
			Token:  "thing-token",
			Name:   "thing",
			Config: configExample,
		}},
		&mocks.FakePublisher{},
	},
	{
		"invalid schema unit",
		"authorization-token",
		"89cf40c23012ce1c",
		[]entities.Config{{
			SensorID: 0,
			Schema: entities.Schema{
				ValueType: 3,
				Unit:      12345,
				TypeID:    65521,
				Name:      "LED",
			},
			Event: entities.Event{
				Change:         true,
				TimeSec:        12,
				LowerThreshold: 25.4,
				UpperThreshold: 87.2,
			},
		}},
		ErrSchemaInvalid,
		false,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{Thing: &entities.Thing{
			ID:     "thing-id",
			Token:  "thing-token",
			Name:   "thing",
			Config: configExample,
		}},
		&mocks.FakePublisher{},
	},
	{
		"invalid schema type ID",
		"authorization-token",
		"89cf40c23012ce1c",
		[]entities.Config{{
			SensorID: 0,
			Schema: entities.Schema{
				ValueType: 3,
				Unit:      0,
				TypeID:    79999,
				Name:      "LED",
			},
			Event: entities.Event{
				Change:         true,
				TimeSec:        12,
				LowerThreshold: 25.4,
				UpperThreshold: 87.2,
			},
		}},
		ErrSchemaInvalid,
		false,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{Thing: &entities.Thing{
			ID:     "thing-id",
			Token:  "thing-token",
			Name:   "thing",
			Config: configExample,
		}},
		&mocks.FakePublisher{},
	},
	{
		"returned information indicating if thing's config has not changed",
		"authorization-token",
		"c09660af89ecba61",
		configExample,
		nil,
		false,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{Thing: &entities.Thing{
			ID:     "thing-id",
			Token:  "thing-token",
			Name:   "thing",
			Config: configExample,
		}, ReturnErr: nil},
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

			thingInteractor := NewThingInteractor(tc.fakeLogger, tc.fakePublisher, tc.fakeThingProxy, &mocks.FakeSessionStore{})
			changed, err := thingInteractor.UpdateConfig(tc.authParam, tc.idParam, tc.configParam)

			assert.EqualValues(t, tc.expectedChanged, changed)
			assert.EqualValues(t, errors.Is(err, tc.expectedError), true)
			tc.fakeThingProxy.AssertExpectations(t)
		})
	}
}
