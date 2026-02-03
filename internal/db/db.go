package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// NewDB creates a new database connection based on the driver type
func NewDB(cfg DBConfig) (*sql.DB, error) {
	if cfg.Driver == "" {
		cfg.Driver = "sqlite3"
	}
	if cfg.Connection == "" {
		cfg.Connection = "./data.db"
	}

	db, err := sql.Open(cfg.Driver, cfg.Connection)
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir banco de dados: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("falha ao conectar ao banco de dados: %w", err)
	}

	// Configurações específicas por driver
	if err := configureDB(db, cfg.Driver); err != nil {
		return nil, err
	}

	if err := createTables(db, cfg.Driver); err != nil {
		return nil, err
	}

	return db, nil
}
