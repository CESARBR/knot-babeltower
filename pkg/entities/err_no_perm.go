package entities

// ErrNoPerm is a error when there is no permission to acess
type ErrNoPerm struct {
	Msg string
}

func (err ErrNoPerm) Error() string {
	return err.Msg
}
