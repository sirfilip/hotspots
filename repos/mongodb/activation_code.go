package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"necsam/config"
	"necsam/errors"
	"necsam/models"
	serializer "necsam/serializers/bson"
)

var (
	ActivationCodeCollection = "activation_codes"
)

// ActivationCodeRepo activation code mongo repo
type ActivationCodeRepo struct {
	client *mongo.Client
}

// Create creates new activation code
func (repo ActivationCodeRepo) Create(ctx context.Context, code models.ActivationCode) error {
	ac := serializer.ActivationCode{}
	if err := ac.Populate(code); err != nil {
		return err
	}
	col := repo.client.Database(config.Get("dbname")).Collection(ActivationCodeCollection)
	_, err := col.InsertOne(ctx, ac)
	return err
}

// FindByCode performs activation code lookup in db by code
func (repo ActivationCodeRepo) FindByCode(ctx context.Context, code string) (models.ActivationCode, error) {
	ac := serializer.ActivationCode{}
	activationCode := models.ActivationCode{}

	coll := repo.client.Database(config.Get("dbname")).Collection(ActivationCodeCollection)
	err := coll.FindOne(ctx, bson.D{{"code", code}}, options.FindOne()).Decode(&ac)
	if err == mongo.ErrNoDocuments {
		return activationCode, errors.RecordNotFound
	}
	return ac.Model(), err
}

// Delete removes existing activation code from db
func (repo ActivationCodeRepo) Delete(ctx context.Context, code models.ActivationCode) error {
	coll := repo.client.Database(config.Get("dbname")).Collection(ActivationCodeCollection)
	_, err := coll.DeleteOne(ctx, bson.D{{"code", code.Code}}, options.Delete())
	return err
}

// NewActivationCodeRepo constructor
func NewActivationCodeRepo(client *mongo.Client) ActivationCodeRepo {
	return ActivationCodeRepo{client}
}
