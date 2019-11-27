package controllers

import (
	"encoding/json"
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

// DetailedErrorResponse represents the response to be sent to the request
type DetailedErrorResponse struct {
	Message string `json:"message"`
}

func (uc *UserController) writeResponse(w http.ResponseWriter, statusCode int, msg interface{}) {
	w.WriteHeader(statusCode)

	if msg == nil {
		return
	}

	js, err := json.Marshal(msg)
	if err != nil {
		uc.logger.Errorf("Unable to marshal json: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		uc.logger.Errorf("Unable to write to connection HTTP: %s", err)
		return
	}
}

func mapErrorToStatusCode(err error) int {
	switch err.(type) {
	case entities.ErrEntityExists:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// Create godoc
// @Summary Creates a new user
// @Produce json
// @Accept  json
// @Param user body entities.User true "User e-mail and password"
// @Success 201 {object} StatusResponse "Message informing the user was created properly"
// @Failure 422 {object} StatusResponse "Invalid request format"
// @Failure 409 {object} StatusResponse "User already exists"
// @Failure 500 {string} string "Internal server error"
// @Router /users [post]
// Create handles the server request and calls CreateUserInteractor
func (uc *UserController) Create(w http.ResponseWriter, r *http.Request) {
	var err error
	var user entities.User

	uc.logger.Debug("Handle request to create user")

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&user)
	if err != nil {
		uc.logger.Error("Failed to parse request body")
		uc.writeResponse(w, http.StatusUnprocessableEntity, nil)
		return
	}

	err = uc.createUserInteractor.Execute(user)
	if err != nil {
		uc.logger.Errorf("Failed to create user")
		der := &DetailedErrorResponse{err.Error()}
		uc.writeResponse(w, mapErrorToStatusCode(err), der)
		return
	}

	uc.logger.Infof("User %s created", user.Email)
	uc.writeResponse(w, http.StatusCreated, nil)
}
