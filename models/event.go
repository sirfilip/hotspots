package models

import (
	"time"

	"github.com/paulmach/orb"

	"necsam/errors"
)

type Event struct {
	EventID     string
	Title       string
	Description string
	Cost        float64
	Date        time.Time
	Location    orb.Geometry
	UserID      string
}

func (e Event) Equals(other Event) bool {
	return e.EventID == other.EventID &&
		e.Title == other.Title &&
		e.Description == other.Description &&
		e.Cost == other.Cost &&
		e.UserID == other.UserID &&
		orb.Equal(e.Location, other.Location)
}

// MarshalJSON prevents marshaling
func (e Event) MarshalJSON() ([]byte, error) {
	return nil, errors.ForbidenMarshalingError
}
