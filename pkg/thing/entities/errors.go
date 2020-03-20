package entities

import "errors"

var (
	// ErrThingForbidden represents the error when thing cannot be authenticated
	ErrThingForbidden = errors.New("forbidden to authenticate thing")

	// ErrThingNotFound represents the error when the schema has a invalid format
	ErrThingNotFound = errors.New("thing not found on thing's service")
)
