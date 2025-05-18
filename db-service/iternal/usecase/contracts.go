package usecase

import (
	"github.com/dami404/diploma-db/iternal/entity"
	"github.com/dami404/diploma-db/iternal/repository"
)

type DBUsecase interface {
	Save(ticket entity.Event) error
	LastEvents(city string) ([]entity.Event, error)
}

type dbUsecase struct {
	repo repository.DBRepository
}

func NewDBUsecase(repo repository.DBRepository) DBUsecase {
	return &dbUsecase{repo: repo}
}
