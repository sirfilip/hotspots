package repos

import (
	"context"

	"necsam/models"
)

type Event interface {
	Create(context.Context, models.Event) error
}
