package entities

import "errors"

var (
	// ErrMalformedRequest represents the error when request has a malformed body or parameters
	ErrMalformedRequest = errors.New("failed due to malformed request")

	// ErrExistingID represents the error when request a key with an already existing ID
	ErrExistingID = errors.New("failed due to using already existing ID")

	// ErrMissingContentType represents the error when request has no content type
	ErrMissingContentType = errors.New("missing or invalid content type")

	// ErrService represents the internal error occured in the service
	ErrService = errors.New("unexpected server-side error occurred")

	// ErrUserForbidden represents the error when user cannot be authenticated
	ErrUserForbidden = errors.New("forbidden to authenticate user")

	// ErrUserExists is returned when trying to register an existing user
	ErrUserExists = errors.New("user is already created")

	// ErrUserBadRequest represents the error when request body is in wrong format
	ErrUserBadRequest = errors.New("unsupported content type, verify e-mail format")
)
