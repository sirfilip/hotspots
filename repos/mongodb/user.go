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
	UserCollection = "users"
)

// UserRepo user repository
type UserRepo struct {
	client *mongo.Client
}

// Create persists user in mongodb
func (repo UserRepo) Create(ctx context.Context, user models.User) error {
	u := serializer.User{}
	if err := u.Populate(user); err != nil {
		return err
	}
	col := repo.client.Database(config.Get("dbname")).Collection(UserCollection)
	_, err := col.InsertOne(ctx, u)
	return err
}

// Update updates existing user in mongodb
func (repo UserRepo) Update(ctx context.Context, user models.User) error {
	u := serializer.User{}
	if err := u.Populate(user); err != nil {
		return err
	}
	col := repo.client.Database(config.Get("dbname")).Collection(UserCollection)
	_, err := col.UpdateOne(ctx, bson.D{{"user_id", u.UserID}}, bson.M{"$set": u})
	return err
}

// Delete removes existing user in mongodb
func (repo UserRepo) Delete(ctx context.Context, user models.User) error {
	coll := repo.client.Database(config.Get("dbname")).Collection(UserCollection)
	_, err := coll.DeleteOne(ctx, bson.D{{"user_id", user.UserID}}, options.Delete())
	return err
}

// FindByID performs user lookup in db by ID
func (repo UserRepo) FindByID(ctx context.Context, userID string) (models.User, error) {
	u := serializer.User{}
	user := models.User{}

	coll := repo.client.Database(config.Get("dbname")).Collection(UserCollection)
	err := coll.FindOne(ctx, bson.D{{"user_id", userID}}, options.FindOne()).Decode(&u)
	if err == mongo.ErrNoDocuments {
		return user, errors.RecordNotFound
	}
	return u.Model(), err
}

// FindByEmail performs user lookup in db by Email
func (repo UserRepo) FindByEmail(ctx context.Context, email string) (models.User, error) {
	u := serializer.User{}
	user := models.User{}

	coll := repo.client.Database(config.Get("dbname")).Collection(UserCollection)
	err := coll.FindOne(ctx, bson.D{{"email", email}}, options.FindOne()).Decode(&u)
	if err == mongo.ErrNoDocuments {
		return user, errors.RecordNotFound
	}
	return u.Model(), err
}

// NewUserRepo UserRepo constructor
func NewUserRepo(client *mongo.Client) UserRepo {
	return UserRepo{client: client}
}
