package entities

// ErrEntityExists is a error when there is conflict with a existent entity
type ErrEntityExists struct {
	Msg string
}

func (err ErrEntityExists) Error() string {
	return err.Msg
}
