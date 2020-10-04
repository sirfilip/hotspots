package bson

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"necsam/models"
)

// ActivationCode bson serializer
type ActivationCode struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    string             `bson:"user_id"`
	Code      string             `bson:"code"`
	CreatedAt time.Time          `bson:"created_at"`
}

// Populate populates user serializer fields from activation code model
func (c *ActivationCode) Populate(activationCode models.ActivationCode) error {
	c.UserID = activationCode.UserID
	c.Code = activationCode.Code
	c.CreatedAt = activationCode.CreatedAt
	return nil
}

// Model creates user model populated from serializer fields
func (c ActivationCode) Model() models.ActivationCode {
	return models.ActivationCode{
		UserID:    c.UserID,
		Code:      c.Code,
		CreatedAt: c.CreatedAt,
	}
}
