package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/lmittmann/tint"
	"github.com/urfave/cli/v3"
)

const Version = "0.1.0"

func main() {
	w := os.Stderr
	slog.SetDefault(slog.New(tint.NewHandler(w, nil)))

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
		slog.Error(fmt.Sprintf("Error: %s", err.Error()))
		os.Exit(1)
	}
}

func getCommand() *cli.Command {
	return &cli.Command{
		Name:  "get",
		Usage: "Fetch information from the EventDoctor API",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "env",
				Aliases: []string{"e"},
				Usage:   "Environment",
				Value:   "development",
			},
			&cli.StringFlag{
				Name:  "topics",
				Usage: "Filter by topics (comma-separated)",
			},
			&cli.StringFlag{
				Name:  "services",
				Usage: "Filter by services (comma-separated)",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			environment := cmd.String("env")
			topics := cmd.String("topics")
			services := cmd.String("services")

			// Validate at least one filter is provided
			if topics == "" && services == "" {
				return fmt.Errorf("at least one filter must be provided (--topics or --services)")
			}

			cmdGet(environment, topics, services)
			return nil
		},
	}
}

// Command implementations (mocked)

func cmdGet(environment string, topics string, services string) {
	fmt.Printf("Getting data with environment: %s\n", environment)
	if topics != "" {
		topicsList := strings.Split(topics, ",")
		fmt.Printf("Topics filter: %v\n", topicsList)
	}
	if services != "" {
		servicesList := strings.Split(services, ",")
		fmt.Printf("Services filter: %v\n", servicesList)
	}
	fmt.Println("[MOCKED] Fetching data from API")
}
