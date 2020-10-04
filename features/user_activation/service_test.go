package user_activation

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"necsam/config"
	"necsam/db"
	"necsam/errors"
	"necsam/models"
	"necsam/repos/mongodb"
)

func TestUserActivcation_Service(t *testing.T) {
	assert := require.New(t)
	ctx := context.TODO()
	client := db.MongoClient()
	defer client.Disconnect(ctx)
	defer client.Database(config.Get("dbname")).Collection(mongodb.UserCollection).Drop(ctx)
	defer client.Database(config.Get("dbname")).Collection(mongodb.ActivationCodeCollection).Drop(ctx)
	user := models.User{UserID: "u-123", Email: "user@example.com", Password: "password", Username: "username"}
	validAC := models.ActivationCode{UserID: user.UserID, Code: "ac-123", CreatedAt: time.Now()}
	expiredAC := models.ActivationCode{UserID: user.UserID, Code: "ac-345", CreatedAt: time.Now().Add(-time.Hour * 25)}
	noUserAC := models.ActivationCode{UserID: "non existing", Code: "ac-567", CreatedAt: time.Now()}
	userRepo := mongodb.NewUserRepo(client)
	acRepo := mongodb.NewActivationCodeRepo(client)

	assert.NoError(userRepo.Create(ctx, user))
	assert.NoError(acRepo.Create(ctx, validAC))
	assert.NoError(acRepo.Create(ctx, expiredAC))
	assert.NoError(acRepo.Create(ctx, noUserAC))

	svc := NewUserActivationService(userRepo, acRepo)

	tests := []struct {
		Title   string
		Code    string
		Error   error
		Finally func()
	}{
		{
			Title: "Performs successful activation for valid activation code",
			Code:  validAC.Code,
			Error: nil,
			Finally: func() {
				user, _ := userRepo.FindByID(ctx, user.UserID)
				assert.Equal(user.Status, models.StatusUserActive)
			},
		},
		{
			Title:   "Responds with RecordNotFound error on code without user",
			Code:    noUserAC.Code,
			Error:   errors.RecordNotFound,
			Finally: func() {},
		},
		{
			Title:   "Responds with RecordNotFound error on non existing code",
			Code:    "non existing code",
			Error:   errors.RecordNotFound,
			Finally: func() {},
		},
		{
			Title:   "Responds with TokenExpiredError on expired token",
			Code:    expiredAC.Code,
			Error:   errors.TokenExpiredError,
			Finally: func() {},
		},
		{
			Title:   "Responds with RecordNotFound on token reuse",
			Code:    validAC.Code,
			Error:   errors.RecordNotFound,
			Finally: func() {},
		},
	}
	for _, test := range tests {
		assert.Equal(test.Error, svc.ActivateUser(ctx, test.Code), test.Title)
		test.Finally()
	}
}
