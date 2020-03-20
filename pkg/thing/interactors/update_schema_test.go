package interactors

import (
	"errors"
	"testing"

	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
	"github.com/CESARBR/knot-babeltower/pkg/thing/mocks"
	"github.com/stretchr/testify/assert"
)

var (
	errThingProxyFailed         = "failed to update the schema on the thing's proxy"
	errPublisherClientFailed    = "failed to send updated schema response"
	errPublisherConnectorFailed = "failed to send update schema to connector"
	errSchemaInvalid            = "invalid schema"
)

type UpdateSchemaTestCase struct {
	name           string
	authorization  string
	thingID        string
	schemaList     []entities.Schema
	isSchemaValid  bool
	expectedErr    error
	fakeLogger     *mocks.FakeLogger
	fakeThingProxy *mocks.FakeThingProxy
	fakePublisher  *mocks.FakePublisher
	fakeConnector  *mocks.FakeConnector
}

var tCases = []UpdateSchemaTestCase{
	{
		"schema successfully updated on the thing's proxy",
		"authorization token",
		"19cf40c23012ce1c",
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
		errors.New(errThingProxyFailed),
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
		&mocks.FakeConnector{},
	},
	{
		"schema response successfully sent",
		"authorization token",
		"39cf40c23012ce1c",
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
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
		&mocks.FakeConnector{},
	},
	{
		"failed to send updated schema response",
		"authorization token",
		"49cf40c23012ce1c",
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
		errors.New(errPublisherClientFailed),
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
		&mocks.FakeConnector{},
	},
	{
		"failed to send update schema to connector",
		"authorization token",
		"59cf40c23012ce1c",
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
		errors.New(errPublisherConnectorFailed),
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
		&mocks.FakeConnector{},
	},
	{
		"invalid schema type ID",
		"authorization token",
		"69cf40c23012ce1c",
		[]entities.Schema{
			{
				SensorID:  0,
				ValueType: 3,
				Unit:      0,
				TypeID:    79999,
				Name:      "LED",
			},
		},
		false,
		nil,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
		&mocks.FakeConnector{},
	},
	{
		"invalid schema unit",
		"authorization token",
		"79cf40c23012ce1c",
		[]entities.Schema{
			{
				SensorID:  0,
				ValueType: 3,
				Unit:      12345,
				TypeID:    65521,
				Name:      "LED",
			},
		},
		false,
		nil,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
		&mocks.FakeConnector{},
	},
	{
		"invalid schema name",
		"authorization token",
		"89cf40c23012ce1c",
		[]entities.Schema{
			{
				SensorID:  0,
				ValueType: 3,
				Unit:      0,
				TypeID:    65521,
				Name:      "SchemaNameGreaterThan23Characters",
			},
		},
		false,
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
				Return(tc.expectedErr).
				Maybe()
			tc.fakeConnector.
				On("SendUpdateSchema", tc.thingID, tc.schemaList).
				Return(tc.expectedErr).
				Maybe()
			tc.fakePublisher.
				On("SendUpdatedSchema", tc.thingID).
				Return(tc.expectedErr).
				Maybe()

			thingInteractor := NewThingInteractor(tc.fakeLogger, tc.fakePublisher, tc.fakeThingProxy, tc.fakeConnector)
			err := thingInteractor.UpdateSchema(tc.authorization, tc.thingID, tc.schemaList)
			if !tc.isSchemaValid {
				assert.EqualError(t, err, errSchemaInvalid)
			}

			tc.fakeThingProxy.AssertExpectations(t)
			tc.fakePublisher.AssertExpectations(t)
		})
	}
}
