package usecase

import (
	"log"

	"github.com/dami404/diploma-parser/iternal/entity"
	"github.com/dami404/diploma-parser/iternal/repository"
)

type parserUsecase struct {
	dbRepo repository.DBRepository
}

func NewParserUsecase(dbRepo repository.DBRepository) ParserUsecase {
	return &parserUsecase{dbRepo: dbRepo}
}

func (uc *parserUsecase) ProfitEvents(event entity.Event) []entity.Event {
	log.Println("Usecase.ParseAndSave")

	// Парсим сайты
	parsedTickets := uc.dbRepo.ProfitEvents(event.City, event.Name)

	return parsedTickets
}
