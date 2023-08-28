package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

const pathMigration = "migrations"

// MustMigrate применяет миграции из переданной директории
func MustMigrate(db *sqlx.DB) {
	if err := goose.Up(db.DB, pathMigration); err != nil {
		panic(fmt.Errorf("apply migrations: %w", err))
	}
}
