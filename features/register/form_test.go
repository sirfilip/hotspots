package register

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"necsam/config"
	"necsam/db"
	"necsam/models"
	"necsam/repos/mongodb"
)

func TestRegisterUser_Validation(t *testing.T) {
	ctx := context.TODO()
	client := db.MongoClient()
	defer client.Disconnect(ctx)
	defer client.Database(config.Get("dbname")).Collection("users").Drop(ctx)
	assert := require.New(t)
	repo := mongodb.NewUserRepo(client)
	user := models.User{UserID: primitive.NewObjectID().Hex(), Email: "email@example.com"}
	assert.NoError(repo.Create(ctx, user))

	tests := []struct {
		Username string
		Email    string
		Password string
		Errors   map[string]string
	}{
		{
			Username: "",
			Email:    "",
			Password: "",
			Errors: map[string]string{
				"username": "Username must be between 3 and 40 chars",
				"email":    "Email is not a valid email",
				"password": "Password must be between 6 and 40 chars",
			},
		},
		{
			Username: "1_132_123",
			Email:    "non valid email",
			Password: "",
			Errors: map[string]string{
				"username": "Username must be alphanumeric",
				"email":    "Email is not a valid email",
				"password": "Password must be between 6 and 40 chars",
			},
		},
	}

	for _, test := range tests {
		form := RegisterForm{
			Username: test.Username,
			Email:    test.Email,
			Password: test.Password,
			repo:     repo,
		}
		errorMessages, err := form.Submit(ctx)
		assert.NoError(err)
		assert.Equal(test.Errors, errorMessages)
	}
}
