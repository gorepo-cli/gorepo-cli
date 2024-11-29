package commands

import (
	"errors"
	"github.com/urfave/cli/v2"
)

func (cmd *Commands) List(c *cli.Context) error {
	if exists := cmd.ConfigManager.RootConfigExists(); !exists {
		return errors.New("monorepo not found at " + cmd.Config.Runtime.ROOT)
	}
	modules, err := cmd.ConfigManager.GetModules()
	if err != nil {
		return err
	}
	if len(modules) == 0 {
		cmd.SystemUtils.Logger.Info("no modules found")
	} else {
		for _, module := range modules {
			cmd.SystemUtils.Logger.Default(module.ModuleConfig.Name)
		}
	}
	return nil
}
