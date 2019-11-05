package controllers

import (
	"net/http"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
)

// UserController represents the controller for user
type UserController struct {
	logger logging.Logger
}

// NewUserController constructs the controller
func NewUserController(logger logging.Logger) *UserController {
	return &UserController{logger}
}

// Create handles the server request
func (uc *UserController) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: parse request
}
