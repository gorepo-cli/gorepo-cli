package commands

import (
	"github.com/urfave/cli/v2"
	"strconv"
)

func (cmd *Commands) Debug(c *cli.Context) error {
	cmd.SystemUtils.Logger.InfoLn("===================")
	cmd.SystemUtils.Logger.InfoLn("RUNTIME_CONFIG")
	cmd.SystemUtils.Logger.InfoLn("===================")
	cmd.SystemUtils.Logger.DefaultLn("WD_(COMMAND_RAN_FROM)........" + cmd.Config.Runtime.WD)
	cmd.SystemUtils.Logger.DefaultLn("ROOT (OF THE MONOREPO)......." + cmd.Config.Runtime.ROOT)
	cmd.SystemUtils.Logger.DefaultLn("MONOREPO EXISTS (AT ROOT)...." +
		strconv.FormatBool(cmd.Config.RootConfigExists()))

	cmd.SystemUtils.Logger.InfoLn("===================")
	cmd.SystemUtils.Logger.InfoLn("STATIC_CONFIG")
	cmd.SystemUtils.Logger.InfoLn("===================")
	cmd.SystemUtils.Logger.DefaultLn("MAX RECURSION................" + strconv.Itoa(cmd.Config.Static.MaxRecursion))
	cmd.SystemUtils.Logger.DefaultLn("ROOT FILE NAME..............." + cmd.Config.Static.RootFileName)
	cmd.SystemUtils.Logger.DefaultLn("MODULE FILE NAME............." + cmd.Config.Static.ModuleFileName)

	if cmd.Config.RootConfigExists() {
		cmd.SystemUtils.Logger.InfoLn("===================")
		cmd.SystemUtils.Logger.InfoLn("ROOT_CONFIG")
		cmd.SystemUtils.Logger.InfoLn("===================")

		cfg, err := cmd.Config.LoadRootConfig()
		if err != nil {
			return err
		}

		cmd.SystemUtils.Logger.DefaultLn("NAME.........." + cfg.Name)
		cmd.SystemUtils.Logger.DefaultLn("VERSION......." + cfg.Version)
		cmd.SystemUtils.Logger.DefaultLn("STRATEGY......" + cfg.Strategy)
		cmd.SystemUtils.Logger.DefaultLn("VENDOR........" + strconv.FormatBool(cfg.Vendor))

		modules, err := cmd.Config.GetModules()
		if err != nil {
			return err
		}

		cmd.SystemUtils.Logger.DefaultLn("N_MODULES....." + strconv.Itoa(len(modules)))

		if len(modules) > 0 {
			cmd.SystemUtils.Logger.InfoLn("===================")
			cmd.SystemUtils.Logger.InfoLn("MODULES_CONFIG")
			cmd.SystemUtils.Logger.InfoLn("===================")
		}

		for _, module := range modules {
			cmd.SystemUtils.Logger.InfoLn("MODULE " + module.Name)
			cmd.SystemUtils.Logger.DefaultLn("MODULE_NAME........ " + module.Name)
			cmd.SystemUtils.Logger.DefaultLn("MODULE_PATH........ " + module.RelativePath)
			if len(module.Scripts) > 0 {
				cmd.SystemUtils.Logger.DefaultLn("COMMANDS........")
				for k, v := range module.Scripts {
					cmd.SystemUtils.Logger.DefaultLn("  " + k + " -> " + v)
				}
			}
		}
	}

	return nil
}
