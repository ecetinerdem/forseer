package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	*sql.DB
}

func NewDB() (*DB, error) {
	dbURL := os.Getenv("DATABASE_URL")

	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	db, err := sql.Open("pgx", dbURL)

	if err != nil {
		return nil, fmt.Errorf("failed to open database %w", err)
	}

	if err := db.PingContext(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return &DB{db}, nil
}
