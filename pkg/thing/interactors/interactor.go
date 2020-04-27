package interactors

import (
	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/thing/delivery/amqp"
	"github.com/CESARBR/knot-babeltower/pkg/thing/delivery/http"
	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
)

// Interactor is an interface that defines the thing's use cases operations
type Interactor interface {
	Register(authorization, id, name string) error
	Unregister(authorization, id string) error
	UpdateSchema(authorization, id string, schemaList []entities.Schema) error
	List(authorization string) ([]*entities.Thing, error)
	RequestData(authorization, thingID string, sensorIds []int) error
	UpdateData(authorization, thingID string, data []entities.Data) error
	PublishData(authorization, thingID string, data []entities.Data) error
	Auth(authorization, id string) error
}

// ThingInteractor represents the thing interactor capabilities, it's composed
// by the necessary dependencies
type ThingInteractor struct {
	logger             logging.Logger
	clientPublisher    amqp.ClientPublisher
	thingProxy         http.ThingProxy
	connectorPublisher amqp.ConnectorPublisher
}

// NewThingInteractor creates a new ThingInteractor instance
func NewThingInteractor(
	logger logging.Logger,
	publisher amqp.ClientPublisher,
	thingProxy http.ThingProxy,
	connector amqp.ConnectorPublisher,
) *ThingInteractor {
	return &ThingInteractor{logger, publisher, thingProxy, connector}
}
