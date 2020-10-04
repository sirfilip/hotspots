package repos

import (
	"context"

	"necsam/models"
)

type ActivationCode interface {
	Create(context.Context, models.ActivationCode) error
	Delete(context.Context, models.ActivationCode) error
	FindByCode(context.Context, string) (models.ActivationCode, error)
}
