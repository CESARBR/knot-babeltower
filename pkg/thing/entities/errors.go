package entities

import "errors"

var (
	// ErrThingUnauthorized represents the error when thing cannot be authenticated
	ErrThingUnauthorized = errors.New("unauthorized to authenticate thing")

	// ErrThingNotFound represents the error when the schema has a invalid format
	ErrThingNotFound = errors.New("thing not found on thing's service")

	// ErrThingExists is returned when trying to register an existing thing
	ErrThingExists = errors.New("thing is already registered")
)
