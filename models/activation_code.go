package models

import (
	"time"

	"necsam/errors"
)

var (
	ExpirationTimeInHours = 24
)

// ActivationCode entity
type ActivationCode struct {
	UserID    string
	Code      string
	CreatedAt time.Time
}

// Equals performs equality test
func (code ActivationCode) Equals(other ActivationCode) bool {
	return code.Code == other.Code && code.UserID == code.UserID
}

// IsExpired checks if the code has been expired
func (code ActivationCode) IsExpired() bool {
	return code.CreatedAt.Add(time.Hour * time.Duration(ExpirationTimeInHours)).Before(time.Now())
}

// MarshalJSON prevents marshaling
func (code ActivationCode) MarshalJSON() ([]byte, error) {
	return nil, errors.ForbidenMarshalingError
}
