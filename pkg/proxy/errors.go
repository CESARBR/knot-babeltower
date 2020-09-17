package proxy

import (
	"fmt"
	"net/http"
)

var (
	// ErrPermissionDenied occurs when request has no permission
	ErrPermissionDenied = fmt.Errorf("invalid credentials or authorization")

	// ErrInvalidThing occurs when indicated thing ID is invalid
	ErrInvalidThing = fmt.Errorf("thing is not registered")

	// ErrExistingEmail occurs when trying to create an user with a registered email
	ErrExistingEmail = fmt.Errorf("email is already registered")
)

// StatusErrors map the response status code to the respective error
var StatusErrors = map[int]error{
	http.StatusForbidden: ErrPermissionDenied,
	http.StatusNotFound:  ErrInvalidThing,
	http.StatusConflict:  ErrExistingEmail,
}
