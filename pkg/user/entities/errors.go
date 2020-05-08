package entities

import "errors"

var (
	// ErrUserForbidden represents the error when user cannot be authenticated
	ErrUserForbidden = errors.New("forbidden to authenticate user")

	// ErrUserExists is returned when trying to register an existing user
	ErrUserExists = errors.New("user is already created")
)
