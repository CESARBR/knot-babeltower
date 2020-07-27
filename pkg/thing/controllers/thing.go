package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/network"
	"github.com/CESARBR/knot-babeltower/pkg/thing/delivery/amqp"
	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
	"github.com/CESARBR/knot-babeltower/pkg/thing/interactors"
)

// ThingController handles the calls to interactor and send to correct topic in queue
type ThingController interface {
	Register(body []byte, authorizationHeader string) error
	Unregister(body []byte, authorizationHeader string) error
	UpdateSchema(body []byte, authorizationHeader string) error
	AuthDevice(body []byte, authorization, replyTo, corrID string) error
	ListDevices(authorization, replyTo, corrID string) error
	PublishData(body []byte, authorization string) error
	RequestData(body []byte, authorization string) error
	UpdateData(body []byte, authorization string) error
}

type thingController struct {
	logger          logging.Logger
	thingInteractor interactors.Interactor
	sender          amqp.Sender
}

// NewThingController constructs the ThingController
func NewThingController(logger logging.Logger, thingInteractor interactors.Interactor, sender amqp.Sender) ThingController {
	return &thingController{logger, thingInteractor, sender}
}

// Register handles the register device request and execute its use case
func (mc *thingController) Register(body []byte, authorizationHeader string) error {
	msg := network.DeviceRegisterRequest{}
	err := json.Unmarshal(body, &msg)
	if err != nil {
		return err
	}

	return mc.thingInteractor.Register(authorizationHeader, msg.ID, msg.Name)
}

// Unregister handles the unregister device request and execute its use case
func (mc *thingController) Unregister(body []byte, authorizationHeader string) error {
	msg := network.DeviceUnregisterRequest{}
	err := json.Unmarshal(body, &msg)
	if err != nil {
		return err
	}

	return mc.thingInteractor.Unregister(authorizationHeader, msg.ID)
}

// UpdateSchema handles the update schema request and execute its use case
func (mc *thingController) UpdateSchema(body []byte, authorizationHeader string) error {
	var updateSchemaReq network.SchemaUpdateRequest
	err := json.Unmarshal(body, &updateSchemaReq)
	if err != nil {
		mc.logger.Error(err)
		return err
	}

	mc.logger.Info("update schema message received")
	mc.logger.Debug(authorizationHeader, updateSchemaReq)

	err = mc.thingInteractor.UpdateSchema(authorizationHeader, updateSchemaReq.ID, updateSchemaReq.Schema)
	if err != nil {
		mc.logger.Error(err)
		return err
	}

	return nil
}

// ListDevices handles the list devices request and execute its use case
func (mc *thingController) ListDevices(authorization, replyTo, corrID string) error {
	mc.logger.Info("list devices command received")
	if replyTo == "" {
		sendErr := mc.sender.SendListResponse([]*entities.Thing{}, replyTo, corrID, interactors.ErrReplyToNotProvided)
		if sendErr != nil {
			return fmt.Errorf("error sending response: %v: %w", interactors.ErrReplyToNotProvided, sendErr)
		}
		return interactors.ErrReplyToNotProvided
	}

	if corrID == "" {
		mc.logger.Warn("Correlation id is empty")
	}

	things, err := mc.thingInteractor.List(authorization)
	if err != nil {
		sendErr := mc.sender.SendListResponse(things, replyTo, corrID, err)
		if sendErr != nil {
			return fmt.Errorf("error sending response: %v: %w", err, sendErr)
		}
		return err
	}

	sendErr := mc.sender.SendListResponse(things, replyTo, corrID, err)
	if sendErr != nil {
		return fmt.Errorf("error sending response: %v: %w", err, sendErr)
	}

	return nil
}

// AuthDevice handles the auth device request and execute its use case
func (mc *thingController) AuthDevice(body []byte, authorization, replyTo, corrID string) error {
	var authThingReq network.DeviceAuthRequest
	err := json.Unmarshal(body, &authThingReq)
	if err != nil {
		mc.logger.Error(err)
		return err
	}
	mc.logger.Info("auth device command received")
	if replyTo == "" {
		sendErr := mc.sender.SendAuthResponse(authThingReq.ID, replyTo, corrID, interactors.ErrReplyToNotProvided)
		if sendErr != nil {
			return fmt.Errorf("error sending response: %v: %w", interactors.ErrReplyToNotProvided, sendErr)
		}
		return interactors.ErrReplyToNotProvided
	}

	if corrID == "" {
		mc.logger.Warn("Correlation id is empty")
	}

	err = mc.thingInteractor.Auth(authorization, authThingReq.ID)
	if err != nil {
		sendErr := mc.sender.SendAuthResponse(authThingReq.ID, replyTo, corrID, err)
		if sendErr != nil {
			return fmt.Errorf("error sending response: %v: %w", err, sendErr)
		}
		return err
	}

	sendErr := mc.sender.SendAuthResponse(authThingReq.ID, replyTo, corrID, err)
	if sendErr != nil {
		return fmt.Errorf("error sending response: %v: %w", err, sendErr)
	}

	return nil
}

// RequestData handles the request data request and execute its use case
func (mc *thingController) RequestData(body []byte, authorization string) error {
	var requestDataReq network.DataRequest
	err := json.Unmarshal(body, &requestDataReq)
	if err != nil {
		mc.logger.Error(err)
		return err
	}

	mc.logger.Info("request data command received")
	mc.logger.Debug(authorization, requestDataReq)
	err = mc.thingInteractor.RequestData(authorization, requestDataReq.ID, requestDataReq.SensorIds)
	if err != nil {
		return err
	}

	return nil
}

// UpdateData handles the update data request and execute its use case
func (mc *thingController) UpdateData(body []byte, authorization string) error {
	msg := network.DataUpdate{}
	err := json.Unmarshal(body, &msg)
	if err != nil {
		return fmt.Errorf("message body parsing error: %w", err)
	}

	return mc.thingInteractor.UpdateData(authorization, msg.ID, msg.Data)
}

// PublishData handles the publish data request and execute its use case
func (mc *thingController) PublishData(body []byte, authorization string) error {
	msg := network.DataSent{}
	err := json.Unmarshal(body, &msg)
	if err != nil {
		return fmt.Errorf("message body parsing error: %w", err)
	}

	return mc.thingInteractor.PublishData(authorization, msg.ID, msg.Data)
}
