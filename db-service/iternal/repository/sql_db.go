package repository

import (
	"context"
	"errors"
	"log"

	"github.com/dami404/diploma-db/iternal/entity"
	"github.com/jackc/pgx/v5"
)

type SQLDBRepository struct {
	db *pgx.Conn
}

func NewDBRepository(ctx context.Context, db *pgx.Conn) *SQLDBRepository {
	return &SQLDBRepository{db: db}
}

func (r *SQLDBRepository) Save(ctx context.Context, ticket entity.Event) error {
	select {
	case <-ctx.Done():
		return errors.New("Repo.Save:context timeout")
	default:
		log.Println("Repo.Save")

		tx, err := r.db.Begin(ctx)
		if err != nil {
			return err
		}

		defer tx.Rollback(ctx)

		// в бд есть триггер, который сам удаляет самую старую запись из таблицы Events, если для 1 города их больше 3
		_, err = tx.Exec(ctx, "INSERT INTO Events (Name, City) VALUES ($1,$2)", ticket.Name, ticket.City)
		if err != nil {
			return err
		}

		return tx.Commit(ctx)
	}
}

func (r *SQLDBRepository) LastEvents(ctx context.Context, city string) ([]entity.Event, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("Repo.LastEvents:context timeout")
	default:
		log.Println("Repo.LastEvents")

		entities := make([]entity.Event, 3)
		var name string
		rows, err := r.db.Query(ctx, "SELECT Name FROM Events WHERE City=$1 ORDER BY id DESC", city)
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			rows.Scan(&name)
			entities = append(entities,
				entity.Event{
					Name: name,
					City: city,
				},
			)
		}

		return entities, nil
	}
}
