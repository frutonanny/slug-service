package database

import (
	"context"
	"database/sql"
)

type txCtxKey struct{}

func NewTxContext(parent context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(parent, txCtxKey{}, tx)
}

func TxFromContext(ctx context.Context) *sql.Tx {
	tx, _ := ctx.Value(txCtxKey{}).(*sql.Tx)
	return tx
}
