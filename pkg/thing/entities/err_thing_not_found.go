package entities

import "fmt"

// ErrThingNotFound represents the error when the schema has a invalid format
type ErrThingNotFound struct {
	ID string
}

func (etnf ErrThingNotFound) Error() string {
	return fmt.Sprintf("Thing %s not found", etnf.ID)
}
