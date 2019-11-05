package controllers

import (
	"net/http"

	"github.com/CESARBR/knot-babeltower/pkg/logging"

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

// Create handles the server request and calls CreateUserInteractor
func (uc *UserController) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: parse request
	uc.logger.Debug("Handle request to create user")

	uc.createUserInteractor.Execute()
}
