package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nicolascb/eventdoctor/internal/client"
	"github.com/nicolascb/eventdoctor/internal/commands"
	"github.com/nicolascb/eventdoctor/internal/logger"
	"github.com/urfave/cli/v3"
)

func configCommand() *cli.Command {
	return &cli.Command{
		Name:  "config",
		Usage: "Manage configuration files",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "env",
				Aliases: []string{"e"},
				Usage:   "Environment",
				Value:   "development",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "validate",
				Usage: "Validate the configuration file",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "file",
						Aliases: []string{"f"},
						Value:   "./eventdoctor.yml",
						Usage:   "Path to configuration file",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					filePath := cmd.String("file")
					return cmdConfigValidate(filePath)
				},
			},
			{
				Name:  "apply",
				Usage: "Apply configuration to the API",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "file",
						Aliases: []string{"f"},
						Value:   "./eventdoctor.yml",
						Usage:   "Path to configuration file",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					environment := cmd.String("env")
					filePath := cmd.String("file")
					return cmdConfigApply(environment, filePath)
				},
			},
		},
	}
}

func cmdConfigValidate(filePath string) error {
	logger.Get().Info(fmt.Sprintf("Validating configuration file: %s", filePath))

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("fail to open config file: %w", err)
	}
	defer file.Close()

	c := commands.NewConfigCommand(nil)
	if err := c.Validate(file); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}
	return nil
}

func cmdConfigApply(environment string, filePath string) error {
	logger.Get().Info(fmt.Sprintf("Applying config file: %s", filePath))

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("fail to open config file: %w", err)
	}
	defer file.Close()

	c := commands.NewConfigCommand(client.NewClient())
	if err := c.Apply(environment, file); err != nil {
		return fmt.Errorf("config apply failed: %w", err)
	}

	return nil
}
