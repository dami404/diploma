package usecase

import (
	"context"
	"log"
	"time"

	"github.com/dami404/diploma-web/internal/entity"
)

func (uc WebUsecase) SaveEvent(event, city string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := uc.repo.SaveEvent(ctx, entity.Event{Name: event, City: city})
	if err != nil {
		log.Fatal(err)
	}
}
