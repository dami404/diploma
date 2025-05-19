package usecase

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/dami404/diploma-parser/iternal/entity"
	"github.com/dami404/diploma-parser/iternal/repository"
)

type parserUsecase struct {
	repo repository.DBRepository
}

func NewParserUsecase(repo repository.DBRepository) ParserUsecase {
	return &parserUsecase{repo: repo}
}

func (uc *parserUsecase) ProfitEvents(event entity.Event) []entity.Event {
	log.Println("Usecase.ProfitEvents")

	// Парсим сайты
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	var mu sync.Mutex
	events := []entity.Event{}
	parsers := []func(context.Context, string, string) []entity.Event{
		// uc.repo.ParseKassir,
		// uc.repo.ParseBileter,
		uc.repo.ParseTicketLand,
	}

	for _, parser := range parsers {
		wg.Add(1)
		go func(parser func(context.Context, string, string) []entity.Event) {
			defer wg.Done()
			parsedEvents := parser(ctx, event.City, event.Name)
			mu.Lock()
			events = append(events, parsedEvents...)
			mu.Unlock()
		}(parser)
	}

	wg.Wait()
	return events
}
