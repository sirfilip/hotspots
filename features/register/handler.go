package register

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Handler(form RegisterForm, svc RegisterService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var validationErrors map[string]string

		ctx := c.Request.Context()

		if err = c.ShouldBindJSON(&form); err == nil {
			if validationErrors, err = form.Submit(ctx); err == nil {
				if len(validationErrors) > 0 {
					c.JSON(http.StatusBadRequest, validationErrors)
					return
				}
				if _, err = svc.Register(ctx, form.Username, form.Email, form.Password); err == nil {
					c.Status(http.StatusCreated)
					return
				}
			}
		}
		log.Printf("ERROR: %v", err)
		c.Status(http.StatusInternalServerError)
	}
}
