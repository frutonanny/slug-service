package database

import (
	"context"
	"database/sql"
	"fmt"
	"io"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

var _ io.Closer = (*DB)(nil)

type DB struct {
	DB    *sqlx.DB
	close func() error
}

func Must(dsn string) *DB {
	db, err := New(dsn)
	if err != nil {
		panic(fmt.Errorf("new database: %v", err))
	}

	return db
}

func New(dsn string) (*DB, error) {
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql open: %v", err)
	}

	return &DB{
		DB: db,
		close: func() error {
			return db.Close()
		},
	}, nil
}

func (db *DB) Close() error {
	if err := db.close(); err != nil {
		return fmt.Errorf("close db: %v", err)
	}

	return nil
}

func (db *DB) RunInTx(ctx context.Context, f func(context.Context) error) (err error) {
	tx := TxFromContext(ctx)
	if tx != nil {
		return f(ctx)
	}

	tx, err = db.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}

		if err = tx.Commit(); err != nil {
			err = fmt.Errorf("commit tx: %w", err)
			return
		}
	}()

	ctx = NewTxContext(ctx, tx)

	return f(ctx)
}

func (db *DB) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return db.getExecutor(ctx).Exec(query, args...)
}

func (db *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return db.getExecutor(ctx).ExecContext(ctx, query, args...)
}

func (db *DB) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return db.getExecutor(ctx).Query(query, args...)
}

func (db *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return db.getExecutor(ctx).QueryContext(ctx, query, args...)
}

func (db *DB) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return db.getExecutor(ctx).QueryRow(query, args...)
}

func (db *DB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return db.getExecutor(ctx).QueryRowContext(ctx, query, args...)
}

type executor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryRow(string, ...interface{}) *sql.Row
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

func (db *DB) getExecutor(ctx context.Context) executor {
	tx := TxFromContext(ctx)
	if tx != nil {
		return tx
	}
	return db.DB
}
