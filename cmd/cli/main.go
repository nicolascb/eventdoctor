package main

import (
	"database/sql"
	"os"

	"go.uber.org/zap"
)

func main() {
}

func handleShutdown(quit chan os.Signal, db *sql.DB, logger *zap.Logger) {
	<-quit
	logger.Info("shutting down server gracefully...")

	if err := db.Close(); err != nil && err != sql.ErrConnDone {
		logger.Error("error closing database", zap.Error(err))
	}

	logger.Info("server stopped successfully")
	logger.Sync()
}
