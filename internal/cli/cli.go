package cli

import (
	"github.com/urfave/cli/v2"
	"gorepo-cli/internal/commands"
	"gorepo-cli/internal/config"
	"gorepo-cli/pkg/systemutils"
	"os"
)

func Exec() (err error) {
	su := systemutils.NewSystemUtils()
	cfg, err := config.NewConfig(su)
	if err != nil {
		return err
	}
	cmd := commands.NewCommands(su, cfg)
	app := &cli.App{
		Name:  "GOREPO",
		Usage: "A CLI tool to manage Go monorepos",
		Commands: []*cli.Command{
			{
				Name:   "init",
				Usage:  "Initialize a new monorepo at the working directory",
				Action: cmd.Init,
			},
			{
				Name:   "list",
				Usage:  "List all modules in the monorepo",
				Action: cmd.List,
			},
			{
				Name:   "run",
				Usage:  "Run a command in a given scope (all modules, some modules, at root)",
				Action: cmd.Run,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "target",
						Value: "all",
						Usage: "Target root or specific modules (comma separated)",
					},
					&cli.BoolFlag{
						Name:  "dry-run",
						Value: false,
						Usage: "Print the commands that would be executed",
					},
					&cli.BoolFlag{
						Name:  "allow-missing",
						Value: false,
						Usage: "Run the scripts in the modules that have it, even if it is missing in some",
					},
				},
			},
			{
				Name:   "version",
				Usage:  "Print the version of the monorepo",
				Action: cmd.Version,
			},
			{
				Name:   "debug",
				Usage:  "Gives information about the configuration",
				Action: cmd.Debug,
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "verbose",
				Usage: "Enable verbose logging for all commands",
				Value: false,
			},
		},
	}
	return app.Run(os.Args)
}
