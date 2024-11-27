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
			},
			{
				Name:   "list",
				Usage:  "List all modules in the monorepo",
				Action: commands.List,
			},
			//{
			//	Name:   "add",
			//	Usage:  "Add a new module to the monorepo",
			//	Action: commands.Add,
			//	Flags: []cli.Flag{
			//		&cli.BoolFlag{
			//			Name:  "verbose",
			//			Usage: "Enable verbose output",
			//		},
			//		&cli.StringFlag{
			//			Name:  "template",
			//			Usage: "Choose a template (not implemented)",
			//		},
			//	},
			//},
			//{
			//	Name:   "run",
			//	Usage:  "Run a command in all modules",
			//	Action: commands.Run,
			//},
			//{}, // sanitize / lint / health / check
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
