package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nicolascb/eventdoctor/internal/logger"
	"github.com/urfave/cli/v3"
)

const Version = "0.1.0"

func main() {
	logger.Init(logger.Config{
		Level:  logger.LevelInfo,
		Format: logger.FormatTint,
		Output: os.Stderr,
	})
	log := logger.Get()

	cmd := &cli.Command{
		Name:    "eventdoctor-cli",
		Version: Version,
		Usage:   "Manage event-driven architecture documentation",
		Commands: []*cli.Command{
			configCommand(),
			getCommand(),
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if len(os.Args) == 1 {
				cli.ShowAppHelp(cmd)
			}
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Error(fmt.Sprintf("Error: %s", err.Error()))
		os.Exit(1)
	}
}
