package bson

import (
	"time"

	"github.com/paulmach/orb"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"necsam/errors"
	"necsam/models"
)

type Point struct {
	Type        string     `bson:"type"`
	Coordinates [2]float64 `bson:"coordinates"`
}

// Event bson serializer for the event model
type Event struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	EventID     string             `bson:"event_id"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
	Cost        float64            `bson:"const"`
	Date        time.Time          `bson:"date"`
	Location    Point              `bson:"location"`
	UserID      string             `bson:"user_id"`
}

// Populate populates event serializer fields from event
func (e *Event) Populate(event models.Event) error {
	e.EventID = event.EventID
	e.Title = event.Title
	e.Description = event.Description
	e.Cost = event.Cost
	e.Date = event.Date
	point, ok := event.Location.(orb.Point)
	if !ok {
		return errors.InvalidGeometryError
	}
	e.Location = Point{
		Type:        point.GeoJSONType(),
		Coordinates: [2]float64{point.Lon(), point.Lat()},
	}
	e.UserID = event.UserID
	return nil
}

// Model creates event model populated from serializer fields
func (e Event) Model() models.Event {
	model := models.Event{
		EventID:     e.EventID,
		Title:       e.Title,
		Description: e.Description,
		Cost:        e.Cost,
		Date:        e.Date,
		UserID:      e.UserID,
	}
	point := orb.Point(e.Location.Coordinates)
	model.Location = point
	return model
}
