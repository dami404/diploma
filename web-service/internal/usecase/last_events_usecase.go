package usecase

import (
	"context"
	"log"
	"time"

	"github.com/dami404/diploma-web/internal/entity"
)

func (uc WebUsecase) LastEvents(city string) []entity.Event {
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()
	events, err := uc.repo.LastEvents(ctx, city)
	if err != nil {
		log.Fatal(err)
	}
	return events
}
