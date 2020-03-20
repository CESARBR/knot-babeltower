package interactors

import "errors"

var (
	// ErrAuthNotProvided is returned when authorization token is not provided
	ErrAuthNotProvided = errors.New("authorization token not provided")

	// ErrIDNotProvided is returned when thing's id is not provided
	ErrIDNotProvided = errors.New("thing's id not provided")

	// ErrDataNotProvided is returned when thing's data is not provided
	ErrDataNotProvided = errors.New("thing's data not provided")
)
