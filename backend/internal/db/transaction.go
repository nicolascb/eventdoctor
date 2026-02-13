package db

import (
	"context"
	"database/sql"
	"fmt"
)

// SQLExecutor abstracts common SQL operations shared by *sql.DB and *sql.Tx.
type SQLExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// WithTransaction runs fn inside a SQL transaction.
// On error the transaction is rolled back; on success it is committed.
func WithTransaction(ctx context.Context, db *sql.DB, fn func(ctx context.Context, tx SQLExecutor) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	err = fn(ctx, tx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction error: %v, rollback failed: %v", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
