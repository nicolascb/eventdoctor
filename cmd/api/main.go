package main

import (
	"fmt"

	"github.com/nicolascb/eventdoctor/internal/api"
	"github.com/nicolascb/eventdoctor/internal/db"
	"go.uber.org/zap"
)

func main() {

	fmt.Println("EventDoctor")
	fmt.Println("Starting EventDoctor API server...")

	cfg, err := api.LoadConfig()
	if err != nil {
		panic(err)
	}

	logger, _ := zap.NewProduction()
	logger.Info("Config loaded", zap.Any("config", cfg))

	store, err := db.NewSQLiteStore(cfg.SQLitePath)
	if err != nil {
		panic(err)
	}

	if err := api.NewAPI(cfg.Port, store, logger).Run(); err != nil {
		panic(err)
	}
}
