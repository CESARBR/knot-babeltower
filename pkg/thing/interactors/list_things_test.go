package interactors

import (
	"errors"
	"testing"

	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
	"github.com/CESARBR/knot-babeltower/pkg/thing/mocks"
	"github.com/stretchr/testify/assert"
)

type listThingsTestCase struct {
	name                        string
	authorization               string
	expectedUseCaseResponse     error
	expectedProxyResponseThings []*entities.Thing
	expectedProxyResponseError  error
	expectedPublisherResponse   error
	fakeLogger                  *mocks.FakeLogger
	fakeThingProxy              *mocks.FakeThingProxy
	fakePublisher               *mocks.FakePublisher
	fakeConnector               *mocks.FakeConnector
}

var ltCases = []listThingsTestCase{
	{
		"authorization key not passed",
		"",
		errors.New("authorization key not provided"),
		[]*entities.Thing{},
		nil,
		nil,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
		&mocks.FakeConnector{},
	},
	{
		"failed to list things from thing's service",
		"authorization token",
		errors.New("Thing's service unavailable"),
		[]*entities.Thing{},
		errors.New("Thing's service unavailable"),
		nil,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
		&mocks.FakeConnector{},
	},
	{
		"things successfully published to message queue",
		"authorization token",
		nil,
		[]*entities.Thing{},
		nil,
		nil,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
		&mocks.FakeConnector{},
	},
	{
		"failed to publish list things response",
		"authorization key",
		errors.New("Message queue unavailable"),
		[]*entities.Thing{},
		nil,
		errors.New("Message queue unavailable"),
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
		&mocks.FakeConnector{},
	},
}

func TestListThings(t *testing.T) {
	for _, tc := range ltCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.fakeThingProxy.
				On("List", tc.authorization).
				Return(tc.expectedProxyResponseThings, tc.expectedProxyResponseError).
				Maybe()
			tc.fakePublisher.
				On("SendDevicesList", tc.expectedProxyResponseThings).
				Return(tc.expectedPublisherResponse).
				Maybe()

			thingInteractor := NewThingInteractor(tc.fakeLogger, tc.fakePublisher, tc.fakeThingProxy, tc.fakeConnector)
			err := thingInteractor.List(tc.authorization)
			if tc.authorization == "" {
				assert.EqualError(t, err, "authorization key not provided")
				return
			}

			if err != nil && !assert.IsType(t, err, tc.expectedUseCaseResponse) {
				t.Errorf("failed to list the devices. Error: %s", err)
				return
			}

			tc.fakeThingProxy.AssertExpectations(t)
			tc.fakePublisher.AssertExpectations(t)
		})
	}
}
