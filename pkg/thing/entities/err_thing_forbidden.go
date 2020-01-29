package entities

// ErrThingForbidden represents the error when thing cannot be authenticated
type ErrThingForbidden struct{}

func (etf ErrThingForbidden) Error() string {
	return "Forbidden to authenticate thing"
}
