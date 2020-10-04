package bson

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"necsam/models"
)

// User bson serializer
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	UserID   string             `bson:"user_id"`
	Username string             `bson:"username"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
	Status   string             `bson:"status"`
}

// Populate populates user serializer fields from user model
func (u *User) Populate(user models.User) error {
	u.UserID = user.UserID
	u.Email = user.Email
	u.Username = user.Username
	u.Password = user.Password
	u.Status = user.Status
	return nil
}

// Model creates user model populated from serializer fields
func (u User) Model() models.User {
	return models.User{
		UserID:   u.UserID,
		Email:    u.Email,
		Username: u.Username,
		Password: u.Password,
		Status:   u.Status,
	}
}
