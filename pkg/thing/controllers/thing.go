package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/network"
	"github.com/CESARBR/knot-babeltower/pkg/thing/interactors"
)

// ThingController handle messages received from a service
type ThingController struct {
	logger          logging.Logger
	thingInteractor interactors.Interactor
}

// NewThingController constructs the ThingController
func NewThingController(logger logging.Logger, thingInteractor interactors.Interactor) *ThingController {
	return &ThingController{logger, thingInteractor}
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

	mc.logger.Info("Update schema message received")
	mc.logger.Debug(authorizationHeader, updateSchemaReq)

	err = mc.thingInteractor.UpdateSchema(authorizationHeader, updateSchemaReq.ID, updateSchemaReq.Schema)
	if err != nil {
		mc.logger.Error(err)
		return err
	}

	return nil
}

// ListDevices handles the list devices request and execute its use case
func (mc *ThingController) ListDevices(authorization string) error {
	return mc.thingInteractor.List(authorization)
}

// AuthDevice handles the auth device request and execute its use case
func (mc *ThingController) AuthDevice(body []byte, authorization string) error {
	var authThingReq network.DeviceAuthRequest
	err := json.Unmarshal(body, &authThingReq)
	if err != nil {
		mc.logger.Error(err)
		return err
	}

	mc.logger.Info("Auth device command received")
	mc.logger.Debug(authorization, authThingReq)
	return mc.thingInteractor.Auth(authorization, authThingReq.ID)
}

// RequestData handles the request data request and execute its use case
func (mc *ThingController) RequestData(body []byte, authorization string) error {
	var requestDataReq network.DataRequest
	err := json.Unmarshal(body, &requestDataReq)
	if err != nil {
		mc.logger.Error(err)
		return err
	}

	mc.logger.Info("Request data command received")
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

	// TODO: call update data interactor
	return nil
}
