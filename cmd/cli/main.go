package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		// Arguments: []cli.Argument{
		// 	&cli.StringArgs{
		// 		Name:      "file",
		// 		UsageText: "path to the configuration file",
		// 	},
		// },
		Commands: commands(),
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func commands() []*cli.Command {
	return []*cli.Command{
		{
			Name:      "config",
			Aliases:   []string{"c"},
			Usage:     "options for config",
			ArgsUsage: "args usage",
			Arguments: []cli.Argument{
				&cli.StringArgs{
					Name:      "file",
					UsageText: "the config file path",
					Min:       0,
					Max:       -1,
				},
			},
			Commands: []*cli.Command{
				{
					Name:    "apply",
					Aliases: []string{"a"},
					Usage:   "Apply EventDoctor config",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Required: true,
							Name:     "api-url",
							Usage:    "The EventDoctor API URL",
							Aliases:  []string{"u"},
						},
					},
					Action: func(ctx context.Context, cmd *cli.Command) error {
						fmt.Println("Apply config", cmd.String("api-url"))
						return nil
					},
				},
			},
		},
	}
}
