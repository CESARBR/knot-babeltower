package entities

import "errors"

var (
	// ErrInvalidTokenType represent the erro when the tokenType provided is not valid
	ErrInvalidTokenType = errors.New("only 'user' and 'app' token types are supported")

	// ErrExistingID represents the error when request a key with an already existing ID
	ErrExistingID = errors.New("failed due to using already existing ID")

	// ErrUserForbidden represents the error when user cannot be authenticated
	ErrUserForbidden = errors.New("forbidden to authenticate user")

	// ErrUserExists is returned when trying to register an existing user
	ErrUserExists = errors.New("user is already created")

	// ErrUserBadRequest represents the error when request body is in wrong format
	ErrUserBadRequest = errors.New("unsupported content type, verify e-mail format")
)
