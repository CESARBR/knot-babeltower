package entities

// ErrInvalidCredentials represents an error when the credentials is invalid
type ErrInvalidCredentials struct {
	Msg string
}

func (err ErrInvalidCredentials) Error() string {
	return err.Msg
}
