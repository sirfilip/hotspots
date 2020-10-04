package repos

import (
	"context"

	"necsam/models"
)

type User interface {
	Create(context.Context, models.User) error
	Update(context.Context, models.User) error
	Delete(context.Context, models.User) error
	FindByID(context.Context, string) (models.User, error)
	FindByEmail(context.Context, string) (models.User, error)
}
