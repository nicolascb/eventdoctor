package main

import (
	"fmt"

	"github.com/nicolascb/eventdoctor/internal/api"
	"go.uber.org/zap"
)

func main() {
	fmt.Println("Starting EventDoctor API server...")

	cfg, err := api.LoadConfig()
	if err != nil {
		panic(err)
	}

	logger, _ := zap.NewProduction()
	logger.Info("Config loaded", zap.Any("config", cfg))
	if err := api.NewAPI(cfg.Port, nil, logger).Run(); err != nil {
		panic(err)
	}
}
