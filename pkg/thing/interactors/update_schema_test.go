package interactors

import (
	"errors"
	"testing"

	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
	"github.com/CESARBR/knot-babeltower/pkg/thing/mocks"
	"github.com/stretchr/testify/assert"
)

var (
	errThingProxyFailed      = errors.New("failed to update the schema on the thing's proxy")
	errPublisherClientFailed = errors.New("failed to send updated schema response")
	errSchemaInvalid         = errors.New("invalid schema")
)

type UpdateSchemaTestCase struct {
	name           string
	authorization  string
	thingID        string
	err            error
	schemaList     []entities.Schema
	isSchemaValid  bool
	expectedErr    error
	fakeLogger     *mocks.FakeLogger
	fakeThingProxy *mocks.FakeThingProxy
	fakePublisher  *mocks.FakePublisher
}

var tCases = []UpdateSchemaTestCase{
	{
		"schema successfully updated on the thing's proxy",
		"authorization-token",
		"19cf40c23012ce1c",
		nil,
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
	},
	{
		"failed to update the schema on the thing's proxy",
		"authorization-token",
		"29cf40c23012ce1c",
		errThingProxyFailed,
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
		errThingProxyFailed,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
	},
	{
		"schema response successfully sent",
		"authorization-token",
		"39cf40c23012ce1c",
		nil,
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
	},
	{
		"failed to send updated schema response",
		"authorization-token",
		"49cf40c23012ce1c",
		errPublisherClientFailed,
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
		errPublisherClientFailed,
		&mocks.FakeLogger{},
		&mocks.FakeThingProxy{},
		&mocks.FakePublisher{},
	},
	{
		"invalid schema type ID",
		"authorization-token",
		"69cf40c23012ce1c",
		errSchemaInvalid,
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
	},
	{
		"invalid schema unit",
		"authorization-token",
		"79cf40c23012ce1c",
		errSchemaInvalid,
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
	},
	{
		"invalid schema name",
		"authorization-token",
		"89cf40c23012ce1c",
		errSchemaInvalid,
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
	},
}

func TestUpdateSchema(t *testing.T) {
	for _, tc := range tCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.fakeThingProxy.
				On("UpdateSchema", tc.thingID, tc.schemaList).
				Return(tc.expectedErr).
				Maybe()

			tc.fakePublisher.
				On("SendUpdatedSchema", tc.thingID, tc.schemaList, tc.err).
				Return(tc.expectedErr).
				Maybe()

			thingInteractor := NewThingInteractor(tc.fakeLogger, tc.fakePublisher, tc.fakeThingProxy)
			err := thingInteractor.UpdateSchema(tc.authorization, tc.thingID, tc.schemaList)
			if !tc.isSchemaValid {
				assert.EqualError(t, err, errSchemaInvalid.Error())
			}

			tc.fakeThingProxy.AssertExpectations(t)
			tc.fakePublisher.AssertExpectations(t)
		})
	}
}
