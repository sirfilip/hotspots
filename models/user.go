package models

import "necsam/errors"

var (
	StatusUserInactive = "inactive"
	StatusUserActive   = "active"
)

// User entity
type User struct {
	UserID    string
	Username  string
	Email     string
	Password  string
	FirstName string
	LastName  string
	Status    string
}

// Equals performs equality test
func (u User) Equals(other User) bool {
	return u.UserID == other.UserID &&
		u.Email == other.Email
}

// MarshalJSON prevents marshaling
func (u User) MarshalJSON() ([]byte, error) {
	return nil, errors.ForbidenMarshalingError
}
