package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
)

const pathMigration = "migrations"

type Database interface {
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

// MustConnect устанавливает соединение с базой.
func MustConnect(dsn string) *sql.DB {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(fmt.Errorf("opening pgx driver: %w", err))
	}

	if err := db.Ping(); err != nil {
		panic(fmt.Errorf("connect db: %w", err))
	}

	return db
}

// MustMigrate применяет миграции из переданной директории
func MustMigrate(db *sql.DB) {
	if err := goose.Up(db, pathMigration); err != nil {
		panic(fmt.Errorf("apply migrations: %w", err))
	}
}
