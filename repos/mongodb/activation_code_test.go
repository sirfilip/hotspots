package mongodb

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"necsam/config"
	"necsam/db"
	"necsam/errors"
	"necsam/models"
	serializer "necsam/serializers/bson"
)

func TestActivationCodeRepo_Create(t *testing.T) {
	ctx := context.TODO()
	client := db.MongoClient()
	defer client.Disconnect(ctx)
	defer client.Database(config.Get("dbname")).Collection(ActivationCodeCollection).Drop(ctx)
	assert := require.New(t)
	repo := NewActivationCodeRepo(client)
	activationCode := models.ActivationCode{Code: "123", CreatedAt: time.Now()}
	ac := serializer.ActivationCode{}
	assert.NoError(repo.Create(ctx, activationCode))
	coll := client.Database(config.Get("dbname")).Collection(ActivationCodeCollection)
	err := coll.FindOne(ctx, bson.D{{"code", activationCode.Code}}).Decode(&ac)
	assert.NoError(err)
	assert.True(activationCode.Equals(ac.Model()))
}

func TestActivationCodeRepo_Delete(t *testing.T) {
	ctx := context.TODO()
	client := db.MongoClient()
	defer client.Disconnect(ctx)
	defer client.Database(config.Get("dbname")).Collection(ActivationCodeCollection).Drop(ctx)
	assert := require.New(t)
	repo := NewActivationCodeRepo(client)
	activationCode := models.ActivationCode{Code: "123", CreatedAt: time.Now()}
	ac := serializer.ActivationCode{}
	assert.NoError(repo.Create(ctx, activationCode))
	assert.NoError(repo.Delete(ctx, activationCode))
	coll := client.Database(config.Get("dbname")).Collection(ActivationCodeCollection)
	err := coll.FindOne(ctx, bson.D{{"code", activationCode.Code}}).Decode(&ac)
	assert.Equal(err, mongo.ErrNoDocuments)
}

func TestActivationCodeRepo_FindByCode(t *testing.T) {
	ctx := context.TODO()
	client := db.MongoClient()
	defer client.Disconnect(ctx)
	defer client.Database(config.Get("dbname")).Collection(ActivationCodeCollection).Drop(ctx)
	assert := require.New(t)
	repo := NewActivationCodeRepo(client)
	activationCode := models.ActivationCode{Code: "123", CreatedAt: time.Now()}
	assert.NoError(repo.Create(ctx, activationCode))

	tests := []struct {
		Code           string
		ActivationCode models.ActivationCode
		Error          error
	}{
		{
			Code:           "not a valid code",
			ActivationCode: models.ActivationCode{},
			Error:          errors.RecordNotFound,
		},
		{
			Code:           activationCode.Code,
			ActivationCode: activationCode,
			Error:          nil,
		},
	}

	for _, test := range tests {
		code, err := repo.FindByCode(ctx, test.Code)
		assert.True(test.ActivationCode.Equals(code))
		assert.Equal(test.Error, err)
	}
}
