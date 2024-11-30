package commands

import "github.com/urfave/cli/v2"

var version = "dev"

func (cmd *Commands) Version(c *cli.Context) error {
	cmd.SystemUtils.Logger.DefaultLn(version)
	return nil
}
