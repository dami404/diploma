package repository

import (
	"github.com/dami404/diploma-parser/iternal/entity"
)

type DBRepository interface {
	ProfitEvents(city string, name string) []entity.Event
}
