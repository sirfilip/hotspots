package json

import (
	"time"

	"github.com/paulmach/orb"

	"necsam/models"
)

// Event json serializer for the event model
type Event struct {
	EventID     string       `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Cost        float64      `json:"const"`
	Date        time.Time    `json:"date"`
	Location    orb.Geometry `json:"location"`
	UserID      string       `json:"user_id"`
}

// Populate populates event serializer fields from event
func (e *Event) Populate(event models.Event) error {
	e.EventID = event.EventID
	e.Title = event.Title
	e.Description = event.Description
	e.Cost = event.Cost
	e.Date = event.Date
	e.Location = event.Location
	e.UserID = event.UserID
	return nil
}

// Model creates event model populated from serializer fields
func (e Event) Model() models.Event {
	return models.Event{
		EventID:     e.EventID,
		Title:       e.Title,
		Description: e.Description,
		Cost:        e.Cost,
		Date:        e.Date,
		Location:    e.Location,
		UserID:      e.UserID,
	}
}
