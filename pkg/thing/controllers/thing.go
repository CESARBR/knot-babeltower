package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/network"
	"github.com/CESARBR/knot-babeltower/pkg/thing/delivery/amqp"
	"github.com/CESARBR/knot-babeltower/pkg/thing/interactors"
)

// ThingController handle messages received from a service
type ThingController struct {
	logger          logging.Logger
	thingInteractor interactors.Interactor
	sender          amqp.Sender
}

// NewThingController constructs the ThingController
func NewThingController(logger logging.Logger, thingInteractor interactors.Interactor, sender amqp.Sender) *ThingController {
	return &ThingController{logger, thingInteractor, sender}
}

// Register handles the register device request and execute its use case
func (mc *ThingController) Register(body []byte, authorizationHeader string) error {
	msg := network.DeviceRegisterRequest{}
	err := json.Unmarshal(body, &msg)
	if err != nil {
		return err
	}

	return mc.thingInteractor.Register(authorizationHeader, msg.ID, msg.Name)
}

// Unregister handles the unregister device request and execute its use case
func (mc *ThingController) Unregister(body []byte, authorizationHeader string) error {
	msg := network.DeviceUnregisterRequest{}
	err := json.Unmarshal(body, &msg)
	if err != nil {
		return err
	}

	return mc.thingInteractor.Unregister(authorizationHeader, msg.ID)
}

// UpdateSchema handles the update schema request and execute its use case
func (mc *ThingController) UpdateSchema(body []byte, authorizationHeader string) error {
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

// UpdateConfig handles the update config request and execute its use case
func (mc *ThingController) UpdateConfig(body []byte, authorizationHeader string) error {
	var updateConfigReq network.ConfigUpdateRequest
	err := json.Unmarshal(body, &updateConfigReq)
	if err != nil {
		mc.logger.Error(err)
		return err
	}

	mc.logger.Info("update config message received")
	mc.logger.Debug(authorizationHeader, updateConfigReq)

	err = mc.thingInteractor.UpdateConfig(authorizationHeader, updateConfigReq.ID, updateConfigReq.Config)
	if err != nil {
		return err
	}

	// TODO: Publish response to message broker

	return nil
}

// ListDevices handles the list devices request and execute its use case
func (mc *ThingController) ListDevices(authorization, replyTo, corrID string) error {
	mc.logger.Info("list devices command received")
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
func (mc *ThingController) AuthDevice(body []byte, authorization, replyTo, corrID string) error {
	var authThingReq network.DeviceAuthRequest
	err := json.Unmarshal(body, &authThingReq)
	if err != nil {
		mc.logger.Error(err)
		return err
	}

	mc.logger.Info("auth device command received")
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
func (mc *ThingController) RequestData(body []byte, authorization string) error {
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
func (mc *ThingController) UpdateData(body []byte, authorization string) error {
	msg := network.DataUpdate{}
	err := json.Unmarshal(body, &msg)
	if err != nil {
		return fmt.Errorf("message body parsing error: %w", err)
	}

	return mc.thingInteractor.UpdateData(authorization, msg.ID, msg.Data)
}

// PublishData handles the publish data request and execute its use case
func (mc *ThingController) PublishData(body []byte, authorization string) error {
	msg := network.DataSent{}
	err := json.Unmarshal(body, &msg)
	if err != nil {
		return fmt.Errorf("message body parsing error: %w", err)
	}

	return mc.thingInteractor.PublishData(authorization, msg.ID, msg.Data)
}
