package services

import (
	"time"

	"github.com/dgrijalva/jwt-go"

	"necsam/config"
	"necsam/models"
)

var (
	TokenLifetimeInMinutes = time.Duration(60)
)

func NewAccessTokenFor(user models.User) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["auth_id"] = user.UserID
	atClaims["username"] = user.Username
	atClaims["exp"] = time.Now().Add(time.Minute * TokenLifetimeInMinutes).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	return at.SignedString([]byte(config.Get("app_secret")))
}
