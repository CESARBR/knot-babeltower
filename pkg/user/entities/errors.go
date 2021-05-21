package entities

import "errors"

var (
	// ErrInvalidTokenType represent the erro when the tokenType provided is not valid
	ErrInvalidTokenType = errors.New("only 'user' and 'app' token types are supported")

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

	// ErrTokenForbidden represents the error when cannot authenticate user or app token
	ErrTokenForbidden = errors.New("failed to authenticate token")

	// ErrUserExists is returned when trying to register an existing user
	ErrUserExists = errors.New("user is already created")

	// ErrUserBadRequest represents the error when request body is in wrong format
	ErrUserBadRequest = errors.New("unsupported content type, verify e-mail format")
)
