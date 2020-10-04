package mongodb

import (
	"context"
	"testing"
	"time"

	"github.com/paulmach/orb"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"

	"necsam/config"
	"necsam/db"
	"necsam/models"
	serializer "necsam/serializers/bson"
)

func TestEventRepo_Create(t *testing.T) {
	assert := require.New(t)
	ctx := context.TODO()
	client := db.MongoClient()
	defer client.Disconnect(ctx)
	defer client.Database(config.Get("dbname")).Collection(EventsCollection).Drop(ctx)
	repo := NewEventRepo(client)
	event := models.Event{
		EventID:     "123",
		Title:       "Going camping",
		Description: "Have great fun camping",
		Cost:        2.34,
		Date:        time.Now().UTC(),
		Location:    orb.Point([2]float64{1.5, 2.3}),
	}
	e := serializer.Event{}
	assert.NoError(repo.Create(ctx, event))

	coll := client.Database(config.Get("dbname")).Collection(EventsCollection)
	err := coll.FindOne(ctx, bson.D{{"event_id", event.EventID}}).Decode(&e)
	assert.NoError(err)
	assert.True(event.Equals(e.Model()))
}
