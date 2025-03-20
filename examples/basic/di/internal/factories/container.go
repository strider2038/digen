package factories

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"basic/di/lookup"
)

func CreateLogger(ctx context.Context, c lookup.Container) *log.Logger {
	return log.New(os.Stdout, "app", log.LstdFlags)
}

func CreateDB(ctx context.Context, c lookup.Container) *sql.DB {
	db, err := sql.Open("postgres", c.Config(ctx).DatabaseURL)
	if err != nil {
		c.SetError(fmt.Errorf("connect to db: %w", err))
	}

	return db
}
