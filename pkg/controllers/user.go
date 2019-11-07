package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CESARBR/knot-babeltower/pkg/logging"

	"github.com/CESARBR/knot-babeltower/pkg/entities"
	"github.com/CESARBR/knot-babeltower/pkg/interactors"
)

// UserController represents the controller for user
type UserController struct {
	logger               logging.Logger
	createUserInteractor *interactors.CreateUser
}

// NewUserController constructs the controller
func NewUserController(
	logger logging.Logger,
	createUserInteractor *interactors.CreateUser) *UserController {
	return &UserController{logger, createUserInteractor}
}

func (uc *UserController) writeResponse(w http.ResponseWriter, status int, err string) {
	js, jsonErr := json.Marshal(StatusResponse{Message: err})
	if jsonErr != nil {
		uc.logger.Errorf("Unable to marshal json: %s", jsonErr)
		return
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	_, writeErr := w.Write(js)
	if writeErr != nil {
		uc.logger.Errorf("Unable to write to connection HTTP: %s", writeErr)
		return
	}
}

func verifyErrorType(err error) int {
	switch err.(type) {
	case entities.ErrEntityExists:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// Create handles the server request and calls CreateUserInteractor
func (uc *UserController) Create(w http.ResponseWriter, r *http.Request) {
	var err error
	var user entities.User
	var decoder *json.Decoder

	uc.logger.Debug("Handle request to create user")

	decoder = json.NewDecoder(r.Body)

	err = decoder.Decode(&user)
	if err != nil {
		errStr := fmt.Sprintf("Invalid request format: %s", err)
		uc.logger.Error(errStr)
		uc.writeResponse(w, http.StatusUnprocessableEntity, errStr)
		return
	}

	err = uc.createUserInteractor.Execute(user)
	if err != nil {
		uc.logger.Errorf("Response error: %s", err)
		uc.writeResponse(w, verifyErrorType(err), err.Error())
		return
	}

	msg := fmt.Sprintf("User %s created", user.Email)
	uc.logger.Info(msg)
	uc.writeResponse(w, http.StatusCreated, msg)
}
