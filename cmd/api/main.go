package main

import (
	"context"
	"database/sql"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nicolascb/eventdoctor/internal/api"
	"github.com/nicolascb/eventdoctor/internal/db"
	"github.com/nicolascb/eventdoctor/internal/service"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	defer logger.Sync()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	srv, db, err := getServer(logger)
	if err != nil {
		logger.Fatal("fail to start server", zap.Error(err))
	}

	go func() {
		if err := srv.Run(); err != nil {
			logger.Error("fail to start http server", zap.Error(err))
		}
	}()

	logger.Info("server started successfully")
	handleShutdown(quit, srv, db, logger)
}

func getServer(logger *zap.Logger) (*api.API, *sql.DB, error) {
	cfg, err := api.LoadConfig()
	if err != nil {
		return nil, nil, err
	}

	database, err := db.NewSQLiteDB(cfg.SQLitePath)
	if err != nil {
		return nil, nil, err
	}

	svc := service.NewService(database)
	apiHTTP := api.NewAPI(cfg.Port, svc, logger)
	return apiHTTP, database, nil
}

func handleShutdown(quit chan os.Signal, srv *api.API, db *sql.DB, logger *zap.Logger) {
	<-quit
	logger.Info("shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("error during server shutdown", zap.Error(err))
	}

	if err := db.Close(); err != nil && err != sql.ErrConnDone {
		logger.Error("error closing database", zap.Error(err))
	}

	logger.Info("server stopped successfully")
	logger.Sync()
}
