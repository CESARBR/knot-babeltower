package entities

// User represents the user domain model
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
}
