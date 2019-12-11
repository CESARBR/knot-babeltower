package interactors

import (
	"errors"
	"testing"

	"github.com/CESARBR/knot-babeltower/pkg/network"
	"github.com/CESARBR/knot-babeltower/pkg/thing/delivery/http"
	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type FakeUpdateSchemaLogger struct{}

type FakeThingProxy struct {
	mock.Mock
}

type FakePublisher struct {
	mock.Mock
}

type FakeThingConnector struct {
	mock.Mock
}

type UpdateSchemaTestCase struct {
	name                   string
	authorization          string
	thingID                string
	schemaList             []entities.Schema
	isSchemaValid          bool
	expectedSchemaResponse error
	expectedUpdatedSchema  error
	fakeLogger             *FakeUpdateSchemaLogger
	fakeThingProxy         *FakeThingProxy
	fakePublisher          *FakePublisher
	fakeConnector          *FakeThingConnector
}

func (fl *FakeUpdateSchemaLogger) Info(...interface{}) {}

func (fl *FakeUpdateSchemaLogger) Infof(string, ...interface{}) {}

func (fl *FakeUpdateSchemaLogger) Debug(...interface{}) {}

func (fl *FakeUpdateSchemaLogger) Warn(...interface{}) {}

func (fl *FakeUpdateSchemaLogger) Error(...interface{}) {}

func (fl *FakeUpdateSchemaLogger) Errorf(string, ...interface{}) {}

func (ftp *FakeThingProxy) Create(id, name, authorization string) (idGenerated string, err error) {
	return "", nil
}

func (ftp *FakeThingProxy) UpdateSchema(authorization, thingID string, schema []entities.Schema) error {
	args := ftp.Called(thingID, schema)
	return args.Error(0)
}

func (ftp *FakeThingProxy) Get(authorization, thingID string) (*http.ThingProxyRepr, error) {
	return nil, nil
}

func (fp *FakePublisher) SendRegisterDevice(network.RegisterResponseMsg) error {
	return nil
}

func (fp *FakePublisher) SendUpdatedSchema(thingID string) error {
	args := fp.Called(thingID)
	return args.Error(0)
}

func (fc *FakeThingConnector) SendRegisterDevice(id, name string) (err error) {
	ret := fc.Called(id, name)
	return ret.Error(0)
}

func (fc *FakeThingConnector) RecvRegisterDevice() (bytes []byte, err error) {
	ret := fc.Called()
	return bytes, ret.Error(1)
}

var cases = []UpdateSchemaTestCase{
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
		&FakeUpdateSchemaLogger{},
		&FakeThingProxy{},
		&FakePublisher{},
		&FakeThingConnector{},
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
		&FakeUpdateSchemaLogger{},
		&FakeThingProxy{},
		&FakePublisher{},
		&FakeThingConnector{},
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
		&FakeUpdateSchemaLogger{},
		&FakeThingProxy{},
		&FakePublisher{},
		&FakeThingConnector{},
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
		errors.New("failed to send updated schema response"),
		&FakeUpdateSchemaLogger{},
		&FakeThingProxy{},
		&FakePublisher{},
		&FakeThingConnector{},
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
		&FakeUpdateSchemaLogger{},
		&FakeThingProxy{},
		&FakePublisher{},
		&FakeThingConnector{},
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
		&FakeUpdateSchemaLogger{},
		&FakeThingProxy{},
		&FakePublisher{},
		&FakeThingConnector{},
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
		&FakeUpdateSchemaLogger{},
		&FakeThingProxy{},
		&FakePublisher{},
		&FakeThingConnector{},
	},
}

func TestUpdateSchema(t *testing.T) {
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.fakeThingProxy.
				On("UpdateSchema", tc.thingID, tc.schemaList).
				Return(tc.expectedSchemaResponse).
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
