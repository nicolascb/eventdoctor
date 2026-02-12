package db

import (
	"context"
	"fmt"
)

// DeleteProducersByRepository remove todos os produtores associados a um repositório
func DeleteProducersByRepository(ctx context.Context, tx SQLExecutor, repository string) (int64, error) {
	query := `DELETE FROM producers WHERE repository = $1`
	result, err := tx.ExecContext(ctx, query, repository)
	if err != nil {
		return 0, fmt.Errorf("erro ao excluir produtores: %w", err)
	}

	return result.RowsAffected()
}

// DeleteConsumersByRepository remove todos os consumidores associados a um repositório
func DeleteConsumersByRepository(ctx context.Context, tx SQLExecutor, repository string) (int64, error) {
	query := `DELETE FROM consumers WHERE repository = $1`
	result, err := tx.ExecContext(ctx, query, repository)
	if err != nil {
		return 0, fmt.Errorf("erro ao excluir consumidores: %w", err)
	}

	return result.RowsAffected()
}
