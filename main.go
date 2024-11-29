package main

import (
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	su := NewSystemUtils()
	cfg, err := NewConfig(su)
	if err != nil {
		su.Logger.Fatal(err.Error())
		os.Exit(1)
	}
	cm := NewConfigManager(su, cfg)
	commands := NewCommands(su, cfg, cm)
	app := &cli.App{
		Name:  "GOREPO",
		Usage: "A CLI tool to manage Go monorepos",
		Commands: []*cli.Command{
			{
				Name:   "init",
				Usage:  "Initialize a new monorepo at the working directory",
				Action: commands.Init,
			},
			{
				Name:   "list",
				Usage:  "List all modules in the monorepo",
				Action: commands.List,
			},
			{
				Name:   "run",
				Usage:  "Run a command in a given scope (all modules, some modules, at root)",
				Action: commands.Run,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "target",
						Usage: "NOT IMPLEMENTED Target root or specific modules (comma separated)",
					},
					&cli.BoolFlag{
						Name:  "dry-run",
						Usage: "NOT IMPLEMENTED Print the commands that would be executed",
					},
				},
			},
			{
				Name:   "debug",
				Usage:  "Gives information about the configuration",
				Action: commands.Debug,
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		su.Logger.Fatal(err.Error())
		os.Exit(1)
	}
}
