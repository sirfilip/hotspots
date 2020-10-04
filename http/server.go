package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"

	"necsam/config"
	"necsam/features/display_user_profile"
	"necsam/features/login"
	"necsam/features/publish_event"
	"necsam/features/register"
	"necsam/features/user_activation"
	"necsam/middleware"
	"necsam/repos/mongodb"
	"necsam/services"
)

type Server struct {
	*gin.Engine
}

func New(client *mongo.Client) *Server {
	srv := &Server{Engine: gin.Default()}

	mailer := services.Mailer{
		EmailHost:         config.Get("email_host"),
		EmailHostUser:     config.Get("email_host_user"),
		EmailHostPassword: config.Get("email_host_password"),
		EmailPort:         config.GetInt("email_port"),
		EmailUseTls:       config.GetBool("email_use_tls"),
	}

	uuidgen := func() (string, error) {
		gen, err := uuid.NewRandom()
		if err != nil {
			return "", err
		}
		return gen.String(), nil
	}
	userRepo := mongodb.NewUserRepo(client)
	acRepo := mongodb.NewActivationCodeRepo(client)
	eventRepo := mongodb.NewEventRepo(client)
	crypter := services.NewBcrypt()
	loginForm := login.NewForm()
	eventForm := publish_event.NewForm()
	loginSvc := login.NewService(userRepo, crypter)
	publishEventSvc := publish_event.NewService(eventRepo, uuidgen)

	srv.Use(middleware.CORSMiddleware)

	health := srv.Group("/health")
	{
		health.GET("/ready", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		health.GET("/live", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})
	}

	auth := srv.Group("/v1/auth")
	{
		auth.POST("/register", register.Handler(
			register.NewForm(userRepo),
			register.NewRegisterService(userRepo, acRepo, crypter, mailer, uuidgen),
		))

		auth.GET("/activate/:code", user_activation.Handler(
			user_activation.NewUserActivationService(userRepo, acRepo),
		))

		auth.POST("/login", login.Handler(loginForm, loginSvc))

	}

	v1 := srv.Group("/v1")
	{
		v1.POST("/events", middleware.AuthRequired, publish_event.Handler(
			eventForm,
			publishEventSvc,
		))

		v1.GET("/users/:id", middleware.AuthRequired, display_user_profile.Handler(
			userRepo,
		))
	}

	return srv
}
