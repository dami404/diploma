package repository

import (
	"context"

	"github.com/dami404/diploma-parser/iternal/entity"
)

type DBRepository interface {
	ParseTicketLand(ctx context.Context, city string, name string) []entity.Event
	ParseKassir(ctx context.Context, city string, name string) []entity.Event
	ParseBileter(ctx context.Context, city string, name string) []entity.Event
}

type HTTPDBRepository struct {
	url string
}

func NewDBRepository(url string) *HTTPDBRepository {
	return &HTTPDBRepository{url: url}
}
