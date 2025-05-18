package repository

import (
	"context"

	"github.com/dami404/diploma-db/iternal/entity"
)

type DBRepository interface {
	Save(ctx context.Context, ticket entity.Event) error
	LastEvents(ctx context.Context, city string) ([]entity.Event, error)
}
