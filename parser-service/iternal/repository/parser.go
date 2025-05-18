package repository

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/dami404/diploma-parser/iternal/entity"
	ws "github.com/dami404/diploma-parser/iternal/repository/websites"
)

type HTTPDBRepository struct {
	url string
}

func NewDBRepository(url string) *HTTPDBRepository {
	return &HTTPDBRepository{url: url}
}

// ProfitEvents fetches events asynchronously using goroutines.
func (r *HTTPDBRepository) ProfitEvents(city string, name string) []entity.Event {
	log.Println("Repository.ProfitEvents")

	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	var mu sync.Mutex
	events := []entity.Event{}
	parsers := []func(context.Context, string, string) []entity.Event{
		ws.ParseKassir,
		// Add other parsers here if needed
	}

	for _, parser := range parsers {
		wg.Add(1)
		go func(parser func(context.Context, string, string) []entity.Event) {
			defer wg.Done()
			parsedEvents := parser(ctx, city, name)
			mu.Lock()
			events = append(events, parsedEvents...)
			mu.Unlock()
		}(parser)
	}

	wg.Wait()
	return events
}
