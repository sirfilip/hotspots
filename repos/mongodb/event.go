package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"necsam/config"
	"necsam/models"
	serializer "necsam/serializers/bson"
)

var (
	EventsCollection = "events"
)

// EventRepo events repository
type EventRepo struct {
	client *mongo.Client
}

// Create persists event in mongodb
func (repo EventRepo) Create(ctx context.Context, event models.Event) error {
	e := serializer.Event{}
	if err := e.Populate(event); err != nil {
		return err
	}

	col := repo.client.Database(config.Get("dbname")).Collection(EventsCollection)
	_, err := col.InsertOne(ctx, e)
	return err
}

func NewEventRepo(client *mongo.Client) EventRepo {
	return EventRepo{client: client}
}
