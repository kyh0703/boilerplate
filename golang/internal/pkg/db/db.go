package db

import (
	"context"
	"database/sql"
	_ "embed"
	"time"

	_ "modernc.org/sqlite"

	"github.com/kyh0703/template/configs"
	"github.com/kyh0703/template/internal/core/domain/model"
)

//go:embed schema.sql
var ddl string

func NewDB(config *configs.Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return nil, err
	}

	return db, nil
}

func NewQueries(db *sql.DB) *model.Queries {
	return model.New(db)
}
