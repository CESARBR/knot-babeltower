package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/CESARBR/knot-babeltower/pkg/logging"

	shared "github.com/CESARBR/knot-babeltower/pkg/entities"
	"github.com/CESARBR/knot-babeltower/pkg/user/entities"
	"github.com/CESARBR/knot-babeltower/pkg/user/interactors"
)

// UserController represents the controller for user
type UserController struct {
	logger                logging.Logger
	createUserInteractor  *interactors.CreateUser
	createTokenInteractor *interactors.CreateToken
}

// NewUserController constructs the controller
func NewUserController(
	logger logging.Logger,
	createUserInteractor *interactors.CreateUser,
	createTokenInteractor *interactors.CreateToken) *UserController {
	return &UserController{logger, createUserInteractor, createTokenInteractor}
}

// CreateTokenResponse is used to map the use case response to HTTP
type CreateTokenResponse struct {
	Token string `json:"token"`
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
	case shared.ErrEntityExists:
		return http.StatusConflict
	case entities.ErrInvalidCredentials:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

// Create godoc
// @Summary Creates a new user
// @Produce json
// @Accept  json
// @Param user body entities.User true "User e-mail and password"
// @Success 201 {object} DetailedErrorResponse "Message informing the user was created properly"
// @Failure 422 {object} DetailedErrorResponse "Invalid request format"
// @Failure 409 {object} DetailedErrorResponse "User already exists"
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

// CreateToken godoc
// @Summary Generate a user's token
// @Produce json
// @Accept  json
// @Param user body entities.User true "User e-mail and password"
// @Success 201 {object} CreateTokenResponse "User's token"
// @Failure 403 {object} DetailedErrorResponse "Invalid credentials"
// @Failure 500 {string} string "Internal server error"
// @Router /tokens [post]
// CreateToken handles the server request and calls CreateTokenInteractor
func (uc *UserController) CreateToken(w http.ResponseWriter, r *http.Request) {
	var user entities.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		uc.logger.Error("Failed to parse request body")
		uc.writeResponse(w, http.StatusUnprocessableEntity, nil)
		return
	}

	token, err := uc.createTokenInteractor.Execute(user)
	if err != nil {
		uc.logger.Errorf("Failed to create user's token: %s", err)
		der := &DetailedErrorResponse{err.Error()}
		uc.writeResponse(w, mapErrorToStatusCode(err), der)
		return
	}

	uc.logger.Infof("User's %s token created", user.Email)
	ctr := &CreateTokenResponse{token}
	uc.writeResponse(w, http.StatusCreated, ctr)
}
