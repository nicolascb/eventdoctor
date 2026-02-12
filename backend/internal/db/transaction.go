package db

import (
	"context"
	"database/sql"
	"fmt"
)

// SQLExecutor é uma interface que abstrai as operações SQL comuns entre *sql.DB e *sql.Tx
type SQLExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// WithTransaction executa uma função dentro de uma transação SQL
// Se a função retornar um erro, a transação é revertida (rollback)
// Se a função for bem-sucedida, a transação é confirmada (commit)
func WithTransaction(ctx context.Context, db *sql.DB, fn func(ctx context.Context, tx SQLExecutor) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}

	err = fn(ctx, tx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("erro na transação: %v, falha ao dar rollback: %v", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("falha ao commitar transação: %w", err)
	}

	return nil
}
