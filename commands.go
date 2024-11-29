package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Commands struct {
	SystemUtils   SystemUtils
	Config        Config
	ConfigManager ConfigManager
}

func NewCommands(su SystemUtils, cfg Config, manager ConfigManager) *Commands {
	return &Commands{
		SystemUtils:   su,
		Config:        cfg,
		ConfigManager: manager,
	}
}

func (cmd *Commands) Init(c *cli.Context) error {
	if exists := cmd.ConfigManager.RootConfigExists(); exists {
		return errors.New("monorepo already exists at " + cmd.Config.Runtime.ROOT)
	}

	config := RootConfig{
		Name:     c.Args().Get(0),
		Version:  "0.1.0",
		Strategy: "workspace",
		Vendor:   true,
	}

	// ask name
	if config.Name == "" {
		reader := bufio.NewReader(os.Stdin)
		base := filepath.Base(cmd.Config.Runtime.ROOT)
		fmt.Print("Enter the monorepo name: " + info(fmt.Sprintf("default: ("+base+")")) + " ")
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		response = strings.TrimSpace(response)
		if response == "" {
			response = base
		}
		config.Name = response
	}

	// ask strategy
	fmt.Println(info("Using go workspace strategy by default (no other option for now)"))

	// ask if should vendor
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Do you want to vendor dependencies? (y/n) " + info("default: (y)") + " ")
	response, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}
	config.Vendor = strings.TrimSpace(response) == "y"

	if config.Strategy == "workspace" {
		if exists := cmd.ConfigManager.GoWorkspaceExists(); !exists {
			err := cmd.SystemUtils.Exec.GoCommand(cmd.Config.Runtime.ROOT, "work", "init")
			if err != nil {
				return err
			}
		}
	} else if config.Strategy == "rewrite" {
		return errors.New("rewrite strategy unsupported yet")
	} else {
		return errors.New("invalid strategy '" + config.Strategy + "'")
	}

	if err := cmd.ConfigManager.WriteRootConfig(config); err != nil {
		return err
	}

	// todo: check existence of modules folder (go.mod) to sanitize everything (create module.toml and make sure they are in the workspace)

	cmd.SystemUtils.Logger.Success("monorepo initialized at " + cmd.Config.Runtime.ROOT)

	return nil
}

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

func (cmd *Commands) Run(c *cli.Context) error {
	if exists := cmd.ConfigManager.RootConfigExists(); !exists {
		return errors.New("monorepo not found at " + cmd.Config.Runtime.ROOT)
	}

	// run all <script_name> or run <script_name>
	// run some <script_name> <module_name> <module_name> <module_name>
	// run root <script_name>

	if c.Args().Len() < 1 {
		return errors.New("no command provided, usage: gorepo run 'command'")
	}

	command := c.Args().Get(0)

	modules, err := cmd.ConfigManager.GetModules()
	if err != nil {
		return err
	}

	// check all modules have the command
	// todo: flag to bypass this check (--allow-missing)
	var modulesWithoutScript []string
	for _, module := range modules {
		if _, ok := module.ModuleConfig.Scripts[command]; !ok {
			modulesWithoutScript = append(modulesWithoutScript, module.ModuleConfig.Name)
		}
	}
	if len(modulesWithoutScript) > 0 {
		return errors.New("following modules are missing the command '" + command + "' :" + strings.Join(modulesWithoutScript, ", "))
	}

	// execute them
	for _, module := range modules {
		cmd.SystemUtils.Logger.Info("running command in " + module.ModuleConfig.Name)
		path := filepath.Join(cmd.Config.Runtime.ROOT, module.RelativePath)
		for _, script := range module.ModuleConfig.Scripts {
			err := cmd.SystemUtils.Exec.BashCommand(path, script)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

//func (cmd *Commands) Add(c *cli.Context) error {
//	//rootConfig, err := cmd.ConfigManager.LoadRootConfig()
//	//if err != nil {
//	//	return err
//	//}
//
//	var path string
//	var name string
//	pathAndName := c.Args().Get(0)
//	if strings.Contains("/", pathAndName) {
//		path = filepath.Dir(pathAndName)
//		name = filepath.Base(pathAndName)
//	} else {
//		name = pathAndName
//	}
//
//	if name == "" {
//		return errors.New("module name is required, use the syntax 'gorepo add module_name'")
//	}
//
//	if exists := cmd.ConfigManager.ModuleFolderExists(name); exists == true {
//		// todo: in the future we could check if it has a module.toml file and create the toml over there
//		return errors.New("module with the name " + name + " already exists")
//	}
//
//	// fmt.Println(rootConfig)
//
//	moduleConfig := ModuleConfig{
//		Name: name,
//	}
//
//	err := cmd.ConfigManager.WriteModuleConfig(moduleConfig, path, name)
//	if err != nil {
//		return err
//	}
//
//	// todo: duplicate module
//
//	return nil
//}

func (cmd *Commands) Debug(c *cli.Context) error {
	cmd.SystemUtils.Logger.Info("===================")
	cmd.SystemUtils.Logger.Info("RUNTIME_CONFIG")
	cmd.SystemUtils.Logger.Info("===================")
	fmt.Println("WD_(COMMAND_RAN_FROM)........" + cmd.Config.Runtime.WD)
	fmt.Println("ROOT (OF THE MONOREPO)......." + cmd.Config.Runtime.ROOT)
	fmt.Println("MONOREPO EXISTS (AT ROOT)...." +
		strconv.FormatBool(cmd.ConfigManager.RootConfigExists()))

	cmd.SystemUtils.Logger.Info("===================")
	cmd.SystemUtils.Logger.Info("STATIC_CONFIG")
	cmd.SystemUtils.Logger.Info("===================")
	fmt.Println("MAX RECURSION................" + strconv.Itoa(cmd.Config.Static.MaxRecursion))
	fmt.Println("ROOT FILE NAME..............." + cmd.Config.Static.RootFileName)
	fmt.Println("MODULE FILE NAME............." + cmd.Config.Static.ModuleFileName)

	if cmd.ConfigManager.RootConfigExists() {
		cmd.SystemUtils.Logger.Info("===================")
		cmd.SystemUtils.Logger.Info("ROOT_CONFIG")
		cmd.SystemUtils.Logger.Info("===================")

		cfg, err := cmd.ConfigManager.LoadRootConfig()
		if err != nil {
			return err
		}

		fmt.Println("NAME.........." + cfg.Name)
		fmt.Println("VERSION......." + cfg.Version)
		fmt.Println("STRATEGY......" + cfg.Strategy)
		fmt.Println("VENDOR........" + strconv.FormatBool(cfg.Vendor))

		modules, err := cmd.ConfigManager.GetModules()
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
			cmd.SystemUtils.Logger.Info("MODULE " + module.ModuleConfig.Name)
			fmt.Println("MODULE_NAME........ " + module.ModuleConfig.Name)
			fmt.Println("MODULE_PATH........ " + module.RelativePath)
			if len(module.ModuleConfig.Scripts) > 0 {
				fmt.Println("COMMANDS........")
				for k, v := range module.ModuleConfig.Scripts {
					fmt.Println("  " + k + " -> " + v)
				}
			}
		}
	}

	return nil
}
