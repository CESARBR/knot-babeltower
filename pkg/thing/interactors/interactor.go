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
	UpdateSchema(authorization, id string, schemaList []entities.Schema) error
}

// ThingInteractor represents the thing interactor capabilities, it's composed
// by the necessary dependencies
type ThingInteractor struct {
	logger       logging.Logger
	msgPublisher amqp.Publisher
	thingProxy   http.ThingProxy
	connector    amqp.Connector
}

// NewThingInteractor creates a new ThingInteractor instance
func NewThingInteractor(
	logger logging.Logger,
	publisher amqp.Publisher,
	thingProxy http.ThingProxy,
	connector amqp.Connector,
) *ThingInteractor {
	return &ThingInteractor{logger, publisher, thingProxy, connector}
}
