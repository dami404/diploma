package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/dami404/diploma-web/internal/entity"
)

type HTTPRepository struct {
	url   string
	dbUrl string
}

func NewHTTPRepository(dbUrl string, url string) *HTTPRepository {
	return &HTTPRepository{
		url:   url,
		dbUrl: dbUrl,
	}
}

func (r *HTTPRepository) ProfitEvents(ctx context.Context, name string, city string) ([]entity.Event, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("Repo.ProfitEvents:context timeout")
	default:
		client := http.Client{}
		url := r.url + "/search/" + city + "/" + name
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		log.Println("Repository.ProfitEvents: response=", strconv.Itoa(resp.StatusCode))
		if resp.StatusCode != http.StatusAccepted {
			return nil, fmt.Errorf("failed to get data: %s", resp.Status)
		}
		var data []entity.Event

		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, fmt.Errorf("Invalid request body", http.StatusBadRequest)
		} else {
			return data, nil
		}
	}
}
func (r *HTTPRepository) LastEvents(ctx context.Context, city string) ([]entity.Event, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("Repo.LastEvents:context timeout")
	default:
		client := http.Client{}
		url := r.dbUrl + "/last/" + city
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, err

		}

		defer resp.Body.Close()
		log.Println("Repository.LastEvents: response=", strconv.Itoa(resp.StatusCode))
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to get data: %s", resp.Status)
		}

		var data []entity.Event

		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, fmt.Errorf("Invalid request body", http.StatusBadRequest)
		} else {
			return data, nil
		}
	}
}

func (r *HTTPRepository) SaveEvent(ctx context.Context, event entity.Event) error {
	select {
	case <-ctx.Done():
		return errors.New("Repo.SaveEvent:context timeout")
	default:
		client := http.Client{}
		url := r.dbUrl + "/save"
		requestBody := map[string]entity.Event{
			"ticket": event,
		}

		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}

		// Создаем POST-запрос
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")

		// Выполняем запрос
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to execute request: %w", err)
		}
		defer resp.Body.Close()

		// Логируем статус ответа
		log.Println("Repository.SaveEvent: response=", strconv.Itoa(resp.StatusCode))

		// Проверяем статус ответа
		if resp.StatusCode != http.StatusCreated {
			return fmt.Errorf("failed to get data: %s", resp.Status)

		}
		return nil
	}
}
