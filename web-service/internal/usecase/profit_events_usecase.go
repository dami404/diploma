package usecase

import (
	"context"
	"log"
	"time"

	"github.com/dami404/diploma-web/internal/entity"
)

func (uc WebUsecase) ProfitEvents(name string, city string) []entity.Event {
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()
	tickets, err := uc.repo.ProfitEvents(ctx, name, city)
	if err != nil {
		log.Fatal(err)
	}
	return tickets

}
