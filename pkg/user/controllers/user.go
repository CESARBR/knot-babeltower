package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/network"

	"github.com/CESARBR/knot-babeltower/pkg/user/entities"
	"github.com/CESARBR/knot-babeltower/pkg/user/interactors"
)

const failedParseRequestBodyMessage = "failed to parse request body"

// UserController represents the controller for user
type UserController struct {
	logger                  logging.Logger
	createUserInteractor    *interactors.CreateUser
	createTokenInteractor   *interactors.CreateToken
	createSessionInteractor *interactors.CreateSession
}

// CreateTokenRequest represents the received parameters for CreateToken operation
type CreateTokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
	Type     string `json:"type"`
	Duration int    `json:"duration"`
}

// CreateTokenResponse is used to map the use case response to HTTP
type CreateTokenResponse struct {
	Token string `json:"token"`
}

// DetailedErrorResponse represents the response to be sent to the request
type DetailedErrorResponse struct {
	Message string `json:"message"`
}

// NewUserController constructs the controller
func NewUserController(
	logger logging.Logger,
	createUserInteractor *interactors.CreateUser,
	createTokenInteractor *interactors.CreateToken,
	createSessionInteractor *interactors.CreateSession) *UserController {
	return &UserController{logger, createUserInteractor, createTokenInteractor, createSessionInteractor}
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

	uc.logger.Debug("handle request to create user")

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&user)
	if err != nil {
		uc.logger.Error(failedParseRequestBodyMessage)
		uc.writeResponse(w, http.StatusUnprocessableEntity, nil)
		return
	}

	err = uc.createUserInteractor.Execute(user)
	if err != nil {
		uc.logger.Errorf("failed to create user")
		der := &DetailedErrorResponse{err.Error()}
		uc.writeResponse(w, mapErrorToStatusCode(err), der)
		return
	}

	uc.logger.Infof("user %s created", user.Email)
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
	var req CreateTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		uc.logger.Error(failedParseRequestBodyMessage)
		uc.writeResponse(w, http.StatusUnprocessableEntity, nil)
		return
	}

	if req.Type == "" {
		req.Type = "user"
	}
	if req.Duration == 0 {
		req.Duration = 60 * 60 * 24 * 365
	}

	user := entities.User{Email: req.Email, Password: req.Password, Token: req.Token}
	token, err := uc.createTokenInteractor.Execute(user, req.Type, req.Duration)
	if err != nil {
		uc.logger.Errorf("failed to create user's token: %s", err)
		der := &DetailedErrorResponse{err.Error()}
		uc.writeResponse(w, mapErrorToStatusCode(err), der)
		return
	}

	uc.logger.Infof("token created for user %s", req.Email)
	ctr := &CreateTokenResponse{token}
	uc.writeResponse(w, http.StatusCreated, ctr)
}

// CreateSession godoc
// @Summary Generate a user's session ID
// @Produce json
// @Accept  json
// @Param user body network.CreateSessionRequest true "User or application token"
// @Success 201 {object} network.CreateSessionResponse "Session ID"
// @Failure 403 {object} DetailedErrorResponse "Invalid credentials"
// @Failure 500 {string} string "Internal server error"
// @Router /sessions [post]
// CreateSession handles the server request and calls CreateSessionInteractor
func (uc *UserController) CreateSession(w http.ResponseWriter, r *http.Request) {
	var req network.CreateSessionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		uc.logger.Error(failedParseRequestBodyMessage)
		uc.writeResponse(w, http.StatusUnprocessableEntity, nil)
		return
	}

	id, err := uc.createSessionInteractor.Execute(req.Token)
	if err != nil {
		uc.logger.Errorf("failed to create user's messaging session: %s", err)
		der := &DetailedErrorResponse{err.Error()}
		uc.writeResponse(w, mapErrorToStatusCode(err), der)
		return
	}

	ctr := &network.CreateSessionResponse{ID: id}
	uc.writeResponse(w, http.StatusCreated, ctr)
}

func (uc *UserController) writeResponse(w http.ResponseWriter, statusCode int, msg interface{}) {
	w.WriteHeader(statusCode)

	if msg == nil {
		return
	}

	js, err := json.Marshal(msg)
	if err != nil {
		uc.logger.Errorf("unable to marshal json: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		uc.logger.Errorf("unable to write to connection HTTP: %s", err)
		return
	}
}

func mapErrorToStatusCode(err error) int {
	switch err {
	case entities.ErrUserForbidden:
		return http.StatusForbidden
	case entities.ErrUserExists:
		return http.StatusConflict
	case entities.ErrUserBadRequest:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
