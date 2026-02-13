package db

import (
	"context"
	"fmt"
)

// DeleteProducersByRepository deletes all producers associated with a repository.
func DeleteProducersByRepository(ctx context.Context, tx SQLExecutor, repository string) (int64, error) {
	query := `DELETE FROM producers WHERE service_id IN (SELECT id FROM services WHERE repository = ?)`
	result, err := tx.ExecContext(ctx, query, repository)
	if err != nil {
		return 0, fmt.Errorf("error deleting producers: %w", err)
	}

	return result.RowsAffected()
}

// DeleteConsumersByRepository deletes all consumers associated with a repository.
func DeleteConsumersByRepository(ctx context.Context, tx SQLExecutor, repository string) (int64, error) {
	query := `DELETE FROM consumers WHERE service_id IN (SELECT id FROM services WHERE repository = ?)`
	result, err := tx.ExecContext(ctx, query, repository)
	if err != nil {
		return 0, fmt.Errorf("error deleting consumers: %w", err)
	}

	return result.RowsAffected()
}
