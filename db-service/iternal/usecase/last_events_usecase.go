package usecase

import (
	"context"
	"time"

	"github.com/dami404/diploma-db/iternal/entity"
)

func (uc *dbUsecase) LastEvents(city string) ([]entity.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()
	return uc.repo.LastEvents(ctx, city)
}
