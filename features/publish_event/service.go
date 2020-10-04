package publish_event

import (
	"context"
	"time"

	"github.com/paulmach/orb"

	"necsam/models"
	"necsam/repos"
)

type Service interface {
	Publish(ctx context.Context, userID, title, description string, cost, latitude, longitude float64, date time.Time) (models.Event, error)
}

type PublishEventService struct {
	eventRepo repos.Event
	uuidgen   func() (string, error)
}

func (svc PublishEventService) Publish(ctx context.Context, userID, title, description string, cost, latitude, longitude float64, date time.Time) (models.Event, error) {
	event := models.Event{
		Title:       title,
		Description: description,
		Cost:        cost,
		Date:        date,
		Location:    orb.Point([2]float64{latitude, longitude}),
		UserID:      userID,
	}

	eventID, err := svc.uuidgen()
	if err != nil {
		return event, err
	}
	event.EventID = "evt-" + eventID

	return event, svc.eventRepo.Create(ctx, event)
}

func NewService(eventRepo repos.Event, uuidgen func() (string, error)) PublishEventService {
	return PublishEventService{eventRepo: eventRepo, uuidgen: uuidgen}
}
