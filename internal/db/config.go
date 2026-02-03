package db

import (
	"database/sql"
	"fmt"
	"time"
)

// DBConfig holds the configuration for database connection
type DBConfig struct {
	Driver     string // "sqlite3", "postgres", "mysql"
	Connection string // DSN or file path
}

func configureDB(db *sql.DB, driver string) error {
	switch driver {
	case "sqlite3":
		// SQLite: ativa foreign keys, WAL mode e busy timeout
		if _, err := db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
			return fmt.Errorf("falha ao ativar foreign_keys: %w", err)
		}
		if _, err := db.Exec(`PRAGMA journal_mode = WAL`); err != nil {
			return fmt.Errorf("falha ao ativar WAL: %w", err)
		}
		if _, err := db.Exec(`PRAGMA busy_timeout = 5000`); err != nil {
			return fmt.Errorf("falha ao configurar busy_timeout: %w", err)
		}

	case "postgres":
		// PostgreSQL: configura pool de conexões para melhor concorrência
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(5 * time.Minute)

	case "mysql":
		// MySQL: configura pool de conexões
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(5 * time.Minute)

	default:
		return fmt.Errorf("driver não suportado: %s", driver)
	}

	return nil
}
