package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nicolascb/eventdoctor/internal/api"
	"github.com/nicolascb/eventdoctor/internal/db"
	"github.com/nicolascb/eventdoctor/internal/logger"
	"github.com/nicolascb/eventdoctor/internal/service"
)

func main() {
	log := logger.Get()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	srv, database, err := getServer(log)
	if err != nil {
		log.Error("fail to start server", slog.Any("error", err))
		os.Exit(1)
	}

	go func() {
		if err := srv.Run(); err != nil {
			log.Error("fail to start http server", slog.Any("error", err))
		}
	}()

	log.Info("server started successfully")
	handleShutdown(quit, srv, database, log)
}

func getServer(log *slog.Logger) (*api.API, *sql.DB, error) {
	cfg, err := api.LoadConfig()
	if err != nil {
		return nil, nil, err
	}

	database, err := db.NewSQLiteDB(cfg.SQLitePath)
	if err != nil {
		return nil, nil, err
	}

	svc := service.NewService(database)
	apiHTTP := api.NewAPI(cfg.Port, svc, log)
	return apiHTTP, database, nil
}

func handleShutdown(quit chan os.Signal, srv *api.API, db *sql.DB, log *slog.Logger) {
	<-quit
	log.Info("shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("error during server shutdown", slog.Any("error", err))
	}

	if err := db.Close(); err != nil && err != sql.ErrConnDone {
		log.Error("error closing database", slog.Any("error", err))
	}

	log.Info("server stopped successfully")
}
