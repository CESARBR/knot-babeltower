package interactors

import (
	"errors"
	"testing"

	"github.com/CESARBR/knot-babeltower/pkg/mocks"
	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
	"github.com/stretchr/testify/assert"
)

type listThingsTestCase struct {
	// descriptive name, e.g. argument not provided
	name string

	// dependencies outputs and operation under test results
	authorization               string
	expectedProxyResponseError  error
	expectedProxyResponseThings []*entities.Thing
	expectedErrorResult         error
	expectedThingsResult        []*entities.Thing

	// mocked dependencies
	fakeLogger     *mocks.FakeLogger
	fakeThingProxy *mocks.FakeThingProxy
}

var things = []*entities.Thing{
	{
		ID:    "8a6f2fe9da74485f",
		Token: "token",
		Name:  "temperature",
	},
}

var ltCases = []listThingsTestCase{
	{
		"authorization token not provided",
		"",
		nil,
		nil,
		ErrAuthNotProvided,
		nil,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
	},
	{
		"failed to list things from thing's service",
		"authorization-token",
		errors.New("thing's service unavailable"),
		nil,
		errors.New("thing's service unavailable"),
		nil,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
	},
	{
		"things successfully received from the thing's service",
		"authorization-token",
		nil,
		things,
		nil,
		things,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
	},
	{
		"return empty list when there is no thing registered on thing's service",
		"authorization-token",
		nil,
		[]*entities.Thing{},
		nil,
		[]*entities.Thing{},
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
	},
}

func TestListThings(t *testing.T) {
	for _, tc := range ltCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.fakeThingProxy.
				On("List", tc.authorization).
				Return(tc.expectedProxyResponseThings, tc.expectedProxyResponseError).
				Maybe()

			thingInteractor := NewThingInteractor(tc.fakeLogger, nil, tc.fakeThingProxy, &mocks.FakeSessionStore{})
			things, err := thingInteractor.List(tc.authorization)
			if tc.authorization == "" {
				assert.EqualError(t, err, ErrAuthNotProvided.Error())
				return
			}

			if err != nil && !errors.As(err, &tc.expectedErrorResult) {
				t.Errorf("failed to list the devices. Error: %s", err)
				return
			}

			if tc.expectedProxyResponseError == nil {
				assert.Equal(t, things, tc.expectedThingsResult)
			}

			tc.fakeThingProxy.AssertExpectations(t)
		})
	}
}
