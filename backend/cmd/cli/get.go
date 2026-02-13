package main

import (
	"context"
	"fmt"

	"github.com/nicolascb/eventdoctor/internal/client"
	"github.com/nicolascb/eventdoctor/internal/commands"
	"github.com/urfave/cli/v3"
)

func getCommand() *cli.Command {
	return &cli.Command{
		Name:  "get",
		Usage: "Fetch information from the EventDoctor API",
		Commands: []*cli.Command{
			getTopicCommand(),
			getServiceCommand(),
		},
	}
}

func serverURLFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "env",
			Aliases: []string{"e"},
			Usage:   "Environment name (used with --file to resolve server URL)",
			Value:   "development",
		},
		&cli.StringFlag{
			Name:    "file",
			Aliases: []string{"f"},
			Usage:   "Path to eventdoctor.yml (used to resolve server URL)",
			Value:   "./eventdoctor.yml",
		},
		&cli.StringFlag{
			Name:    "url",
			Aliases: []string{"u"},
			Usage:   "Direct server URL (overrides --file and --env)",
		},
	}
}

func getTopicCommand() *cli.Command {
	return &cli.Command{
		Name:      "topic",
		Usage:     "Show producers and consumers for a topic",
		ArgsUsage: "<topic-name>",
		Flags:     serverURLFlags(),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() == 0 {
				return fmt.Errorf("topic name is required")
			}

			c := commands.NewGetCommand(client.NewClient())

			serverURL, err := c.ResolveServerURL(cmd.String("file"), cmd.String("env"), cmd.String("url"))
			if err != nil {
				return err
			}

			return c.Topic(serverURL, cmd.Args().First())
		},
	}
}

func getServiceCommand() *cli.Command {
	return &cli.Command{
		Name:      "service",
		Usage:     "Show producers and consumers for a service",
		ArgsUsage: "<service-name>",
		Flags:     serverURLFlags(),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() == 0 {
				return fmt.Errorf("service name is required")
			}

			c := commands.NewGetCommand(client.NewClient())

			serverURL, err := c.ResolveServerURL(cmd.String("file"), cmd.String("env"), cmd.String("url"))
			if err != nil {
				return err
			}

			return c.Service(serverURL, cmd.Args().First())
		},
	}
}
