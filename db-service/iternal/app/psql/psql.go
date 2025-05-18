package psql

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func NewConnection(ctx context.Context, url string) *pgx.Conn {
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return conn
}
