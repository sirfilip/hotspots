package user_activation

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"necsam/errors"
)

// Handler http interface for user activation
func Handler(svc UserActivator) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error

		code := c.Param("code")
		if err = svc.ActivateUser(c.Request.Context(), code); err == nil {
			c.Status(http.StatusNoContent)
			return
		}

		if err == errors.RecordNotFound {
			c.Status(http.StatusNotFound)
			return
		}

		if err == errors.TokenExpiredError {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Code is expired",
			})
			return
		}

		log.Printf("ERROR: got error: %v", err)
		c.Status(http.StatusInternalServerError)
	}
}
