package interactors

import (
	"errors"
	"testing"

	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
	"github.com/CESARBR/knot-babeltower/pkg/thing/mocks"
	"github.com/stretchr/testify/assert"
)

type UpdateSchemaTestCase struct {
	name                          string
	authorization                 string
	thingID                       string
	schemaList                    []entities.Schema
	isSchemaValid                 bool
	expectedSchemaResponse        error
	expectedUpdatedSchema         error
	expectedUpdateSchemaConnector error
	fakeLogger                    *mocks.FakeLogger
	fakeThingProxy                *mocks.FakeThingProxy
	fakePublisher                 *mocks.FakePublisher
	fakeConnector                 *mocks.FakeConnector
}

var tCases = []UpdateSchemaTestCase{
	{
		"schema successfully updated on the thing's proxy",
		"authorization token",
		"29cf40c23012ce1c",
		[]entities.Schema{
			{
				SensorID:  0,
				ValueType: 3,
				Unit:      0,
				TypeID:    65521,
				Name:      "LED",
			},
		},
		true,
		nil,
		nil,
		nil,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
		&mocks.FakeConnector{},
	},
	{
		"failed to update the schema on the thing's proxy",
		"authorization token",
		"29cf40c23012ce1c",
		[]entities.Schema{
			{
				SensorID:  0,
				ValueType: 3,
				Unit:      0,
				TypeID:    65521,
				Name:      "LED",
			},
		},
		true,
		errors.New("failed to update schema"),
		nil,
		nil,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
		&mocks.FakeConnector{},
	},
	{
		"schema response successfully sent",
		"authorization token",
		"29cf40c23012ce1c",
		[]entities.Schema{
			{
				SensorID:  0,
				ValueType: 3,
				Unit:      0,
				TypeID:    65521,
				Name:      "LED",
			},
		},
		true,
		nil,
		nil,
		nil,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
		&mocks.FakeConnector{},
	},
	{
		"failed to send updated schema response",
		"authorization token",
		"29cf40c23012ce1c",
		[]entities.Schema{
			{
				SensorID:  0,
				ValueType: 3,
				Unit:      0,
				TypeID:    65521,
				Name:      "LED",
			},
		},
		true,
		nil,
		nil,
		errors.New("failed to send updated schema response"),
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
		&mocks.FakeConnector{},
	},
	{
		"failed to send update schema to connector",
		"authorization token",
		"29cf40c23012ce1c",
		[]entities.Schema{
			{
				SensorID:  0,
				ValueType: 3,
				Unit:      0,
				TypeID:    65521,
				Name:      "LED",
			},
		},
		true,
		nil,
		nil,
		errors.New("failed to send update schema to connector"),
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
		&mocks.FakeConnector{},
	},
	{
		"invalid schema type ID",
		"authorization token",
		"29cf40c23012ce1c",
		[]entities.Schema{
			{
				SensorID:  0, // invalid ID
				ValueType: 3,
				Unit:      0,
				TypeID:    79999,
				Name:      "LED",
			},
		},
		false,
		nil,
		nil,
		nil,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
		&mocks.FakeConnector{},
	},
	{
		"invalid schema unit",
		"authorization token",
		"29cf40c23012ce1c",
		[]entities.Schema{
			{
				SensorID:  0, // invalid ID
				ValueType: 3,
				Unit:      12345,
				TypeID:    65521,
				Name:      "LED",
			},
		},
		false,
		nil,
		nil,
		nil,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
		&mocks.FakeConnector{},
	},
	{
		"invalid schema name",
		"authorization token",
		"29cf40c23012ce1c",
		[]entities.Schema{
			{
				SensorID:  0, // invalid ID
				ValueType: 3,
				Unit:      0,
				TypeID:    65521,
				Name:      "SchemaNameGreaterThan23Characters",
			},
		},
		false,
		nil,
		nil,
		nil,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
		&mocks.FakeConnector{},
	},
}

func TestUpdateSchema(t *testing.T) {
	for _, tc := range tCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.fakeThingProxy.
				On("UpdateSchema", tc.thingID, tc.schemaList).
				Return(tc.expectedSchemaResponse).
				Maybe()
			tc.fakeConnector.
				On("SendUpdateSchema", tc.thingID, tc.schemaList).
				Return(tc.expectedUpdateSchemaConnector).
				Maybe()
			tc.fakePublisher.
				On("SendUpdatedSchema", tc.thingID).
				Return(tc.expectedUpdatedSchema).
				Maybe()

			thingInteractor := NewThingInteractor(tc.fakeLogger, tc.fakePublisher, tc.fakeThingProxy, tc.fakeConnector)
			err := thingInteractor.UpdateSchema(tc.authorization, tc.thingID, tc.schemaList)
			if !tc.isSchemaValid {
				assert.EqualError(t, err, "Thing's schema is invalid")
			}

			tc.fakeThingProxy.AssertExpectations(t)
			tc.fakePublisher.AssertExpectations(t)
		})
	}
}
