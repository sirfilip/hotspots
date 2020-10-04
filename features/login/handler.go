package login

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"necsam/errors"
	"necsam/models"
	serializer "necsam/serializers/json"
)

func Handler(form Form, svc Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var validationErrors map[string]string
		var token models.Token

		ctx := c.Request.Context()

		if err = c.ShouldBindJSON(&form); err == nil {
			if validationErrors, err = form.Submit(ctx); err == nil {
				if len(validationErrors) > 0 {
					c.JSON(http.StatusBadRequest, gin.H{
						"message": "Wrong username and password combination",
					})
					return
				}
				if token, err = svc.Login(ctx, form.Email, form.Password); err == nil {
					jsonSerializer := serializer.Token{}
					jsonSerializer.Populate(token)
					c.JSON(http.StatusOK, jsonSerializer)
					return
				}
			}
		}
		if err == errors.RecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Wrong username and password combination",
			})
			return
		}

		log.Printf("ERROR: %v", err)
		c.Status(http.StatusInternalServerError)
	}
}
