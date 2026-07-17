package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type DBTX interface {
	Rebind(query string) string
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error)
	GetContext(ctx context.Context, dest any, query string, args ...any) error
	SelectContext(ctx context.Context, dest any, query string, args ...any) error
}

type contextKey string

const txKey contextKey = "tx"

func (r *SQLiteRepository) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	if _, ok := ctx.Value(txKey).(DBTX); ok {
		return fn(ctx)
	}

	tx, err := r.DB.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	newCtx := context.WithValue(ctx, txKey, tx)

	if err := fn(newCtx); err != nil {
		return fmt.Errorf("fail in the transaction: %v", err)
	}

	return tx.Commit()
}

func (r *SQLiteRepository) getTx(ctx context.Context) DBTX {
	if tx, ok := ctx.Value(txKey).(DBTX); ok {
		return tx
	}
	return r.DB
}
