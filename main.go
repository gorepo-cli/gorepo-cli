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
		Name:  "GOREPO-CLI",
		Usage: "A CLI tool to manage Go monorepos",
		Commands: []*cli.Command{
			{
				Name:   "init",
				Usage:  "Initialize a new monorepo at the working directory",
				Action: commands.Init,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "verbose",
						Usage: "Enable verbose output",
					},
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		su.Logger.Fatal(err.Error())
		os.Exit(1)
	}
}
