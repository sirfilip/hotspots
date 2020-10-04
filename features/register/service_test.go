package register

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"necsam/config"
	"necsam/db"
	"necsam/repos/mongodb"
	"necsam/services"
)

func TestRegisterService_RegisterSuccess(t *testing.T) {
	ctx := context.TODO()
	assert := require.New(t)
	client := db.MongoClient()
	defer client.Disconnect(ctx)
	defer client.Database(config.Get("dbname")).Collection(mongodb.UserCollection).Drop(ctx)
	defer client.Database(config.Get("dbname")).Collection(mongodb.ActivationCodeCollection).Drop(ctx)
	userRepo := mongodb.NewUserRepo(client)
	acRepo := mongodb.NewActivationCodeRepo(client)
	uuidgen := func() (string, error) {
		gen, err := uuid.NewRandom()
		if err != nil {
			return "", err
		}
		return gen.String(), nil
	}
	mailer := services.Mailer{
		EmailHost:         config.Get("email_host"),
		EmailHostUser:     config.Get("email_host_user"),
		EmailHostPassword: config.Get("email_host_password"),
		EmailPort:         config.GetInt("email_port"),
		EmailUseTls:       config.GetBool("email_use_tls"),
	}
	crypter := services.NewBcrypt()
	svc := NewRegisterService(userRepo, acRepo, crypter, mailer, uuidgen)
	_, err := svc.Register(ctx, "username", "email@example.com", "password")
	assert.NoError(err)
}
