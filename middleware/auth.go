package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"necsam/config"
)

func AuthRequired(c *gin.Context) {
	tokenString := strings.Replace(c.GetHeader("Authorization"), "Bearer ", "", -1)
	if tokenString == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Get("app_secret")), nil
	})
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if _, ok = claims["authorized"]; !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("auth_id", claims["auth_id"])
	c.Next()
}
