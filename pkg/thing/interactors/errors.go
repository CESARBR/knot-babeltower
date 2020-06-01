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

	// ErrIDLength is returned when the thing's id have more than 16 ascii characters
	ErrIDLength = errors.New("id length exceeds 16 characters")

	// ErrIDNotHex is returned when the thing's id is not formatted in hexadecimal base
	ErrIDNotHex = errors.New("id is not in hexadecimal format")

	// ErrSchemaInvalid is returned when schema has an invalid format
	ErrSchemaInvalid = errors.New("invalid schema")

	// ErrSensorInvalid is returned when some sensorId mismatch with thing's schema
	ErrSensorInvalid = errors.New("sensor list is incompatible with thing's schema")

	// ErrSchemaUndefined is returned when the thing has no schema yet
	ErrSchemaUndefined = errors.New("thing has no schema")

	// ErrDataInvalid is returned when the provided data mismatch the thing's schema
	ErrDataInvalid = errors.New("data is incompatible with thing's schema")

	// ErrReplyToNotProvided is returned when the reply_to is not provided in RPC calls
	ErrReplyToNotProvided = errors.New("reply_to property not provided")
)
