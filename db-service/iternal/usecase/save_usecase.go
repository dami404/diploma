package usecase

import (
	"context"
	"time"

	"github.com/dami404/diploma-db/iternal/entity"
)

func (uc *dbUsecase) Save(ticket entity.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()
	return uc.repo.Save(ctx, ticket)
}
