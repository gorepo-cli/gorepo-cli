package commands

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"strconv"
)

func (cmd *Commands) Debug(c *cli.Context) error {
	cmd.SystemUtils.Logger.Info("===================")
	cmd.SystemUtils.Logger.Info("RUNTIME_CONFIG")
	cmd.SystemUtils.Logger.Info("===================")
	fmt.Println("WD_(COMMAND_RAN_FROM)........" + cmd.Config.Runtime.WD)
	fmt.Println("ROOT (OF THE MONOREPO)......." + cmd.Config.Runtime.ROOT)
	fmt.Println("MONOREPO EXISTS (AT ROOT)...." +
		strconv.FormatBool(cmd.Config.RootConfigExists()))

	cmd.SystemUtils.Logger.Info("===================")
	cmd.SystemUtils.Logger.Info("STATIC_CONFIG")
	cmd.SystemUtils.Logger.Info("===================")
	fmt.Println("MAX RECURSION................" + strconv.Itoa(cmd.Config.Static.MaxRecursion))
	fmt.Println("ROOT FILE NAME..............." + cmd.Config.Static.RootFileName)
	fmt.Println("MODULE FILE NAME............." + cmd.Config.Static.ModuleFileName)

	if cmd.Config.RootConfigExists() {
		cmd.SystemUtils.Logger.Info("===================")
		cmd.SystemUtils.Logger.Info("ROOT_CONFIG")
		cmd.SystemUtils.Logger.Info("===================")

		cfg, err := cmd.Config.LoadRootConfig()
		if err != nil {
			return err
		}

		fmt.Println("NAME.........." + cfg.Name)
		fmt.Println("VERSION......." + cfg.Version)
		fmt.Println("STRATEGY......" + cfg.Strategy)
		fmt.Println("VENDOR........" + strconv.FormatBool(cfg.Vendor))

		modules, err := cmd.Config.GetModules()
		if err != nil {
			return err
		}

		fmt.Println("N_MODULES....." + strconv.Itoa(len(modules)))

		if len(modules) > 0 {
			cmd.SystemUtils.Logger.Info("===================")
			cmd.SystemUtils.Logger.Info("MODULES_CONFIG")
			cmd.SystemUtils.Logger.Info("===================")
		}

		for _, module := range modules {
			cmd.SystemUtils.Logger.Info("MODULE " + module.Name)
			fmt.Println("MODULE_NAME........ " + module.Name)
			fmt.Println("MODULE_PATH........ " + module.RelativePath)
			if len(module.Scripts) > 0 {
				fmt.Println("COMMANDS........")
				for k, v := range module.Scripts {
					fmt.Println("  " + k + " -> " + v)
				}
			}
		}
	}

	return nil
}
