package usecase

import (
	"github.com/dami404/diploma-parser/iternal/entity"
)

type ParserUsecase interface {
	ProfitEvents(event entity.Event) []entity.Event
}
