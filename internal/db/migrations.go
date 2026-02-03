package db

import (
	"database/sql"
	"embed"
	"fmt"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func createTables(db *sql.DB, driver string) error {
	// Determina qual arquivo de migration usar
	var migrationFile string
	switch driver {
	case "sqlite3":
		migrationFile = "migrations/migration.sql"
	case "postgres":
		migrationFile = "migrations/migration.postgres.sql"
	case "mysql":
		migrationFile = "migrations/migration.mysql.sql"
	default:
		return fmt.Errorf("driver não suportado: %s", driver)
	}

	migrationBytes, err := migrationsFS.ReadFile(migrationFile)
	if err != nil {
		return fmt.Errorf("falha ao ler migration (%s): %w", migrationFile, err)
	}

	_, err = db.Exec(string(migrationBytes))
	if err != nil {
		return fmt.Errorf("falha ao executar migrations: %w", err)
	}

	return nil
}
