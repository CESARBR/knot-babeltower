package interactors

import (
	"errors"
	"testing"

	"github.com/CESARBR/knot-babeltower/pkg/mocks"
	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
	"github.com/stretchr/testify/assert"
)

type GetDataTestCase struct {
	name                        string
	authorization               string
	thingID                     string
	sensorIds                   []int
	expectedThing               *entities.Thing
	expectedThingError          error
	expectedRequestDataResponse error
	fakeLogger                  *mocks.FakeLogger
	fakeThingProxy              *mocks.FakeThingProxy
	fakePublisher               *mocks.FakePublisher
}

var gdCases = []GetDataTestCase{
	{
		"authorization token not provided",
		"",
		"",
		nil,
		nil,
		nil,
		nil,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
	},
	{
		"failed to authenticate with provided token",
		"authorization-token",
		"fc3fcf912d0c290a",
		nil,
		nil,
		errors.New("invalid credentials"),
		nil,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
	},
	{
		"thing doesn't exists on thing's service",
		"authorization-token",
		"fc3fcf912d0c290a",
		nil,
		nil,
		errors.New("thing fc3fcf912d0c290a not found"),
		nil,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
	},
	{
		"thing successfully obtained from the thing's service",
		"authorization-token",
		"fc3fcf912d0c290a",
		[]int{2},
		&entities.Thing{
			ID:    "fc3fcf912d0c290a",
			Token: "token",
			Name:  "thing",
			Config: []entities.Config{
				{
					SensorID: 0,
					Schema: entities.Schema{
						ValueType: 2,
						Unit:      0,
						TypeID:    65521,
						Name:      "test",
					},
					Event: entities.Event{
						Change:         true,
						TimeSec:        12,
						LowerThreshold: 25.4,
						UpperThreshold: 87.2,
					},
				},
			},
		},
		nil,
		nil,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
	},
	{
		"thing hasn't schema for the requested sensor",
		"authorization-token",
		"fc3fcf912d0c290a",
		[]int{1}, // the sensor id 1 can't be mapped to thing's schema
		&entities.Thing{
			ID:    "fc3fcf912d0c290a",
			Token: "token",
			Name:  "thing",
			Config: []entities.Config{
				{
					SensorID: 0,
					Schema: entities.Schema{
						ValueType: 2,
						Unit:      0,
						TypeID:    65521,
						Name:      "test",
					},
					Event: entities.Event{
						Change:         true,
						TimeSec:        12,
						LowerThreshold: 25.4,
						UpperThreshold: 87.2,
					},
				},
			},
		},
		nil,
		nil,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
	},
	{
		"failed to send quest data command to message queue",
		"authorization-token",
		"fc3fcf912d0c290a",
		[]int{1},
		&entities.Thing{
			ID:    "fc3fcf912d0c290a",
			Token: "token",
			Name:  "thing",
			Config: []entities.Config{
				{
					SensorID: 0,
					Schema: entities.Schema{
						ValueType: 2,
						Unit:      0,
						TypeID:    65521,
						Name:      "test",
					},
					Event: entities.Event{
						Change:         true,
						TimeSec:        12,
						LowerThreshold: 25.4,
						UpperThreshold: 87.2,
					},
				},
			},
		},
		nil,
		errors.New("failed to send request data message"),
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
	},
	{
		"request data command successfully sent",
		"authorization-token",
		"fc3fcf912d0c290a",
		[]int{1},
		&entities.Thing{
			ID:    "fc3fcf912d0c290a",
			Token: "token",
			Name:  "thing",
			Config: []entities.Config{
				{
					SensorID: 0,
					Schema: entities.Schema{
						ValueType: 2,
						Unit:      0,
						TypeID:    65521,
						Name:      "test",
					},
					Event: entities.Event{
						Change:         true,
						TimeSec:        12,
						LowerThreshold: 25.4,
						UpperThreshold: 87.2,
					},
				},
			},
		},
		nil,
		nil,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
	},
}

func TestGetData(t *testing.T) {
	for _, tc := range gdCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.fakeThingProxy.
				On("Get", tc.authorization, tc.thingID).
				Return(tc.expectedThing, tc.expectedThingError).
				Maybe()
			tc.fakePublisher.
				On("PublishRequestData", tc.thingID, tc.sensorIds).
				Return(tc.expectedRequestDataResponse).
				Maybe()
		})

		thingInteractor := NewThingInteractor(tc.fakeLogger, tc.fakePublisher, tc.fakeThingProxy)
		err := thingInteractor.RequestData(tc.authorization, tc.thingID, tc.sensorIds)
		if tc.authorization == "" {
			assert.EqualError(t, err, ErrAuthNotProvided.Error())
		}

		tc.fakeThingProxy.AssertExpectations(t)
		tc.fakePublisher.AssertExpectations(t)
	}
}
