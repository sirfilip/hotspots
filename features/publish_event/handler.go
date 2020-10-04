package publish_event

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"necsam/models"
	serializer "necsam/serializers/json"
)

func Handler(form Form, svc Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var validationErrors map[string]string
		var err error
		var event models.Event

		ctx := c.Request.Context()

		authID, ok := c.Get("auth_id")
		if !ok {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if err = c.ShouldBindJSON(&form); err == nil {
			if validationErrors, err = form.Submit(ctx); err == nil {
				if len(validationErrors) > 0 {
					c.JSON(http.StatusBadRequest, validationErrors)
					return
				}

				if event, err = svc.Publish(
					ctx,
					authID.(string),
					form.Title,
					form.Description,
					form.Cost,
					form.Latitude,
					form.Longitude,
					form.EventTime,
				); err == nil {
					jsonSerializer := serializer.Event{}
					jsonSerializer.Populate(event)
					c.JSON(http.StatusCreated, jsonSerializer)
					return
				}
			}
		}

		log.Printf("ERROR: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}
