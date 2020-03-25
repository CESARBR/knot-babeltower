package interactors

import "errors"

var (
	// ErrAuthNotProvided is returned when authorization token is not provided
	ErrAuthNotProvided = errors.New("authorization token not provided")

	// ErrIDNotProvided is returned when thing's id is not provided
	ErrIDNotProvided = errors.New("thing's id not provided")

	// ErrNameNotProvided is returned when thing's name is not provided
	ErrNameNotProvided = errors.New("thing's name not provided")

	// ErrSchemaNotProvided is returned when thing's schema is not provided
	ErrSchemaNotProvided = errors.New("thing's schema not provided")

	// ErrDataNotProvided is returned when thing's data is not provided
	ErrDataNotProvided = errors.New("thing's data not provided")

	// ErrSensorsNotProvided is returned when thing's sensors are not provided
	ErrSensorsNotProvided = errors.New("thing's sensors not provided")
)
