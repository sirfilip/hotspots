package display_user_profile

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"necsam/errors"
	"necsam/models"
	"necsam/repos"
	serializer "necsam/serializers/json"
)

// Handler presents http interface for accessing client
func Handler(repo repos.User) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var user models.User

		userID := c.Param("id")

		if user, err = repo.FindByID(c.Request.Context(), userID); err == nil {
			serializr := serializer.User{}
			serializr.Populate(user)
			c.JSON(http.StatusOK, serializr)
			return
		}

		if errors.Is(errors.RecordNotFound, err) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		log.Printf("ERROR: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}
