package mongodb

import (
	"context"
	"necsam/config"
	"necsam/db"
	"necsam/errors"
	"necsam/models"
	serializer "necsam/serializers/bson"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestUserRepo_Create(t *testing.T) {
	ctx := context.TODO()
	client := db.MongoClient()
	defer client.Disconnect(ctx)
	defer client.Database(config.Get("dbname")).Collection(UserCollection).Drop(ctx)
	assert := require.New(t)
	repo := NewUserRepo(client)
	user := models.User{UserID: "123"}
	u := serializer.User{}
	assert.NoError(repo.Create(ctx, user))
	coll := client.Database(config.Get("dbname")).Collection(UserCollection)
	err := coll.FindOne(ctx, bson.D{{"user_id", user.UserID}}).Decode(&u)
	assert.NoError(err)
	assert.True(user.Equals(u.Model()))
}

func TestUserRepo_Update(t *testing.T) {
	ctx := context.TODO()
	client := db.MongoClient()
	defer client.Disconnect(ctx)
	defer client.Database(config.Get("dbname")).Collection(UserCollection).Drop(ctx)
	assert := require.New(t)
	repo := NewUserRepo(client)
	user := models.User{UserID: "123"}
	u := serializer.User{}
	assert.NoError(repo.Create(ctx, user))
	user.Email = "updated@example.com"
	assert.NoError(repo.Update(ctx, user))
	coll := client.Database(config.Get("dbname")).Collection(UserCollection)
	err := coll.FindOne(ctx, bson.D{{"user_id", user.UserID}}).Decode(&u)
	assert.NoError(err)
	assert.True(user.Equals(u.Model()))
}

func TestUserRepo_Delete(t *testing.T) {
	ctx := context.TODO()
	client := db.MongoClient()
	defer client.Disconnect(ctx)
	defer client.Database(config.Get("dbname")).Collection(UserCollection).Drop(ctx)
	assert := require.New(t)
	repo := NewUserRepo(client)
	user := models.User{UserID: "123", Email: "email@example.com"}
	u := serializer.User{}
	assert.NoError(repo.Create(ctx, user))
	assert.NoError(repo.Delete(ctx, user))
	coll := client.Database(config.Get("dbname")).Collection(UserCollection)
	err := coll.FindOne(ctx, bson.D{{"user_id", user.UserID}}).Decode(&u)
	assert.Equal(err, mongo.ErrNoDocuments)
}

func TestUserRepo_FindByID(t *testing.T) {
	ctx := context.TODO()
	client := db.MongoClient()
	defer client.Disconnect(ctx)
	defer client.Database(config.Get("dbname")).Collection(UserCollection).Drop(ctx)
	assert := require.New(t)
	repo := NewUserRepo(client)
	user := models.User{UserID: "123", Email: "email@example.com"}
	assert.NoError(repo.Create(ctx, user))

	tests := []struct {
		ID    string
		User  models.User
		Error error
	}{
		{
			ID:    "not a valid user id",
			User:  models.User{},
			Error: errors.RecordNotFound,
		},
		{
			ID:    user.UserID,
			User:  user,
			Error: nil,
		},
	}

	for _, test := range tests {
		user, err := repo.FindByID(ctx, test.ID)
		assert.Equal(test.User, user)
		assert.Equal(test.Error, err)
	}
}

func TestUserRepo_FindByEmail(t *testing.T) {
	ctx := context.TODO()
	client := db.MongoClient()
	defer client.Disconnect(ctx)
	defer client.Database(config.Get("dbname")).Collection(UserCollection).Drop(ctx)
	assert := require.New(t)
	repo := NewUserRepo(client)
	user := models.User{UserID: "123", Email: "email@example.com"}
	assert.NoError(repo.Create(ctx, user))

	tests := []struct {
		Email string
		User  models.User
		Error error
	}{
		{
			Email: "Non existing",
			User:  models.User{},
			Error: errors.RecordNotFound,
		},
		{
			Email: user.Email,
			User:  user,
			Error: nil,
		},
	}

	for _, test := range tests {
		user, err := repo.FindByEmail(ctx, test.Email)
		assert.Equal(test.User, user)
		assert.Equal(test.Error, err)
	}
}
