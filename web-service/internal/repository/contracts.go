package repository

import (
	"context"

	"github.com/dami404/diploma-web/internal/entity"
)

type Repository interface {
	ProfitEvents(ctx context.Context, name string, city string) ([]entity.Event, error)
	LastEvents(ctx context.Context, city string) ([]entity.Event, error)
	SaveEvent(ctx context.Context, event entity.Event) error
}
