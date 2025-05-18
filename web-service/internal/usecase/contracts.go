package usecase

import (
	"github.com/dami404/diploma-web/internal/entity"
	"github.com/dami404/diploma-web/internal/repository"
)

type Usecase interface {
	LastEvents(city string) []entity.Event
	ProfitEvents(name string, city string) []entity.Event
	SaveEvent(event, city string)
}

type WebUsecase struct {
	repo repository.Repository
}

func NewUsecase(
	repo repository.Repository) *WebUsecase {
	return &WebUsecase{
		repo: repo,
	}
}
