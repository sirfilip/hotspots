package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"

	"necsam/config"
	"necsam/db"
	"necsam/features/publish_event"
	"necsam/models"
	"necsam/repos/mongodb"
	"necsam/services"
)

func TestServer_Register(t *testing.T) {
	ctx := context.TODO()
	client := db.MongoClient()
	defer client.Disconnect(ctx)
	defer client.Database(config.Get("dbname")).Collection(mongodb.UserCollection).Drop(ctx)
	defer client.Database(config.Get("dbname")).Collection(mongodb.ActivationCodeCollection).Drop(ctx)
	assert := require.New(t)
	server := New(client)
	ts := httptest.NewServer(server)
	defer ts.Close()

	apiClient := resty.New()

	tests := []struct {
		Payload  map[string]interface{}
		Response string
		Status   int
	}{
		{
			Payload: map[string]interface{}{
				"username": "tester",
				"email":    "tester@example.com",
				"password": "password",
			},
			Response: ``,
			Status:   201,
		},
		{
			Payload: map[string]interface{}{
				"username": "",
				"email":    "",
				"password": "",
			},
			Response: `{
					"username": "Username must be between 3 and 40 chars",
					"email":    "Email is not a valid email",
					"password": "Password must be between 6 and 40 chars"
			}`,
			Status: 400,
		},
		{
			Payload: map[string]interface{}{
				"username": "1_132_123",
				"email":    "invalid email",
				"password": "",
			},
			Response: `{
					"username": "Username must be alphanumeric",
					"email":    "Email is not a valid email",
					"password": "Password must be between 6 and 40 chars"
			}`,
			Status: 400,
		},
	}

	for _, test := range tests {
		resp, err := apiClient.R().
			SetHeader("Content-Type", "application/json").
			SetBody(test.Payload).
			Post(ts.URL + "/v1/auth/register")
		assert.NoError(err)
		assert.Equal(test.Status, resp.StatusCode())
		if test.Status == 400 {
			assert.JSONEq(test.Response, resp.String())
		}
	}
}

func TestServer_UserActivation(t *testing.T) {
	ctx := context.TODO()
	client := db.MongoClient()
	defer client.Disconnect(ctx)
	defer client.Database(config.Get("dbname")).Collection(mongodb.UserCollection).Drop(ctx)
	defer client.Database(config.Get("dbname")).Collection(mongodb.ActivationCodeCollection).Drop(ctx)
	assert := require.New(t)
	server := New(client)
	ts := httptest.NewServer(server)
	defer ts.Close()

	userRepo := mongodb.NewUserRepo(client)
	acRepo := mongodb.NewActivationCodeRepo(client)

	user := models.User{UserID: "u-123", Email: "user@example.com", Password: "password", Username: "username"}
	validAC := models.ActivationCode{UserID: user.UserID, Code: "ac-123", CreatedAt: time.Now()}
	expiredAC := models.ActivationCode{UserID: user.UserID, Code: "ac-345", CreatedAt: time.Now().Add(-time.Hour * 25)}
	noUserAC := models.ActivationCode{UserID: "non existing", Code: "ac-567", CreatedAt: time.Now()}

	assert.NoError(userRepo.Create(ctx, user))
	assert.NoError(acRepo.Create(ctx, validAC))
	assert.NoError(acRepo.Create(ctx, expiredAC))
	assert.NoError(acRepo.Create(ctx, noUserAC))

	apiClient := resty.New()

	tests := []struct {
		Title      string
		Code       string
		StatusCode int
	}{
		{
			Title:      "Performs successful activation for valid activation code",
			Code:       validAC.Code,
			StatusCode: http.StatusNoContent,
		},
		{
			Title:      "Responds with RecordNotFound error on code without user",
			Code:       noUserAC.Code,
			StatusCode: http.StatusNotFound,
		},
		{
			Title:      "Responds with RecordNotFound error on non existing code",
			Code:       "non-existing-code",
			StatusCode: http.StatusNotFound,
		},
		{
			Title:      "Responds with TokenExpiredError on expired token",
			Code:       expiredAC.Code,
			StatusCode: http.StatusBadRequest,
		},
		{
			Title:      "Responds with RecordNotFound on token reuse",
			Code:       validAC.Code,
			StatusCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		resp, err := apiClient.R().
			SetHeader("Content-Type", "application/json").
			Get(ts.URL + "/v1/auth/activate/" + test.Code)
		assert.NoError(err)
		assert.Equal(test.StatusCode, resp.StatusCode())
	}
}

func TestServer_UserDisplay(t *testing.T) {
	ctx := context.TODO()
	client := db.MongoClient()
	defer client.Disconnect(ctx)
	defer client.Database(config.Get("dbname")).Collection(mongodb.UserCollection).Drop(ctx)
	assert := require.New(t)
	server := New(client)
	ts := httptest.NewServer(server)
	defer ts.Close()

	userRepo := mongodb.NewUserRepo(client)

	user := models.User{UserID: "u-123", Email: "user@example.com", Password: "password", Username: "username"}
	token, err := services.NewAccessTokenFor(user)
	assert.NoError(err)

	assert.NoError(userRepo.Create(ctx, user))

	apiClient := resty.New()

	tests := []struct {
		Title      string
		ID         string
		StatusCode int
	}{
		{
			Title:      "Performs successfull user display",
			ID:         user.UserID,
			StatusCode: http.StatusOK,
		},
		{
			Title:      "Responds with RecordNotFound error on non existing user id",
			ID:         "an invalid user id",
			StatusCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		resp, err := apiClient.R().
			SetHeader("Authorization", "Bearer "+token).
			SetHeader("Content-Type", "application/json").
			Get(ts.URL + "/v1/users/" + test.ID)
		assert.NoError(err)
		assert.Equal(test.StatusCode, resp.StatusCode())
	}
}

func TestServer_Login(t *testing.T) {
	ctx := context.TODO()
	client := db.MongoClient()
	defer client.Disconnect(ctx)
	defer client.Database(config.Get("dbname")).Collection(mongodb.UserCollection).Drop(ctx)
	defer client.Database(config.Get("dbname")).Collection(mongodb.ActivationCodeCollection).Drop(ctx)
	assert := require.New(t)
	server := New(client)
	ts := httptest.NewServer(server)
	defer ts.Close()

	userRepo := mongodb.NewUserRepo(client)
	activeUser := models.User{Username: "username", Email: "username@example.com", Password: "password"}
	inactiveUser := models.User{Username: "inactive", Email: "inactive@example.com", Password: "password"}

	apiClient := resty.New()

	for _, user := range []models.User{activeUser, inactiveUser} {
		resp, err := apiClient.R().
			SetHeader("Content-Type", "application/json").
			SetBody(map[string]interface{}{
				"username": user.Username,
				"email":    user.Email,
				"password": user.Password,
			}).
			Post(ts.URL + "/v1/auth/register")
		assert.NoError(err)
		assert.Equal(http.StatusCreated, resp.StatusCode())
	}

	user, err := userRepo.FindByEmail(ctx, activeUser.Email)
	assert.NoError(err)
	user.Status = models.StatusUserActive
	assert.NoError(userRepo.Update(ctx, user))

	tests := []struct {
		Title      string
		Email      string
		Password   string
		StatusCode int
	}{
		{
			Title:      "Activated user",
			Email:      activeUser.Email,
			Password:   activeUser.Password,
			StatusCode: http.StatusOK,
		},
		{
			Title:      "Inactive user",
			Email:      inactiveUser.Email,
			Password:   inactiveUser.Password,
			StatusCode: http.StatusBadRequest,
		},
		{
			Title:      "Wrong username and password",
			Email:      "wrongemail@example.com",
			Password:   "wrongpassword",
			StatusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		resp, err := apiClient.R().
			SetHeader("Content-Type", "application/json").
			SetBody(map[string]interface{}{
				"email":    test.Email,
				"password": test.Password,
			}).
			Post(ts.URL + "/v1/auth/login")
		assert.NoError(err)
		assert.Equal(test.StatusCode, resp.StatusCode())
	}
}

func TestServer_PublishEvent(t *testing.T) {
	ctx := context.TODO()
	client := db.MongoClient()
	defer client.Disconnect(ctx)
	defer client.Database(config.Get("dbname")).Collection(mongodb.EventsCollection).Drop(ctx)
	assert := require.New(t)
	server := New(client)
	ts := httptest.NewServer(server)
	defer ts.Close()

	apiClient := resty.New()

	user := models.User{UserID: "123", Email: "user@example.com", Username: "username"}
	accessToken, err := services.NewAccessTokenFor(user)
	assert.NoError(err)

	tests := []struct {
		TestTitle   string
		Title       string
		Description string
		Cost        float64
		Date        string
		Latitude    float64
		Longitude   float64
		Token       string
		StatusCode  int
	}{
		{
			TestTitle:   "Valid Event",
			Title:       "going shopping",
			Description: "Have a great time",
			Cost:        5.6,
			Date:        time.Now().Format(publish_event.DateFormat),
			Latitude:    80.0,
			Longitude:   80.0,
			Token:       accessToken,
			StatusCode:  http.StatusCreated,
		},
		{
			TestTitle:   "Bad token",
			Title:       "going shopping",
			Description: "Have a great time",
			Cost:        5.6,
			Date:        time.Now().Format(publish_event.DateFormat),
			Latitude:    80.0,
			Longitude:   80.0,
			Token:       "wrong access token",
			StatusCode:  http.StatusUnauthorized,
		},
		{
			TestTitle:   "Invalid Event: Missing title",
			Title:       "",
			Description: "Have a great time",
			Cost:        5.6,
			Date:        time.Now().Format(publish_event.DateFormat),
			Latitude:    80.0,
			Longitude:   80.0,
			Token:       accessToken,
			StatusCode:  http.StatusBadRequest,
		},
		{
			TestTitle:   "Invalid Event: Missing Description",
			Title:       "A title",
			Description: "",
			Cost:        5.6,
			Date:        time.Now().Format(publish_event.DateFormat),
			Latitude:    80.0,
			Longitude:   80.0,
			Token:       accessToken,
			StatusCode:  http.StatusBadRequest,
		},
		{
			TestTitle:   "Invalid Event: Wrong Date Format",
			Title:       "A title",
			Description: "Description",
			Cost:        5.6,
			Date:        "2020-31-02",
			Latitude:    80.0,
			Longitude:   80.0,
			Token:       accessToken,
			StatusCode:  http.StatusBadRequest,
		},
		{
			TestTitle:   "Invalid Event - Invalid Latitude",
			Title:       "going shopping",
			Description: "Have a great time",
			Cost:        5.6,
			Date:        time.Now().Format(publish_event.DateFormat),
			Latitude:    500.0,
			Longitude:   80.0,
			Token:       accessToken,
			StatusCode:  http.StatusBadRequest,
		},
		{
			TestTitle:   "Invalid Event - Invalid Longitude",
			Title:       "going shopping",
			Description: "Have a great time",
			Cost:        5.6,
			Date:        time.Now().Format(publish_event.DateFormat),
			Latitude:    80.0,
			Longitude:   500.0,
			Token:       accessToken,
			StatusCode:  http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		resp, err := apiClient.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Authorization", "Bearer "+test.Token).
			SetBody(map[string]interface{}{
				"title":       test.Title,
				"description": test.Description,
				"cost":        test.Cost,
				"date":        test.Date,
				"latitude":    test.Latitude,
				"longitude":   test.Longitude,
			}).
			Post(ts.URL + "/v1/events")
		assert.NoError(err)
		assert.Equal(test.StatusCode, resp.StatusCode(), test.TestTitle)
	}
}
