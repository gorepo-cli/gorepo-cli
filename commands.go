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

	rootConfig := RootConfig{
		Name:     c.Args().Get(0),
		Version:  "0.1.0",
		Strategy: "workspace",
		Vendor:   true,
	}

	// ask name
	if rootConfig.Name == "" {
		reader := bufio.NewReader(os.Stdin)
		defaultName := filepath.Base(cmd.Config.Runtime.ROOT)
		fmt.Print("Enter the monorepo name: " + info(fmt.Sprintf("default: ("+defaultName+")")) + " ")
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		response = strings.TrimSpace(response)
		if response == "" {
			response = defaultName
		}
		rootConfig.Name = response
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
	rootConfig.Vendor = strings.TrimSpace(response) == "y"

	// handle go workspace
	if rootConfig.Strategy == "workspace" {
		if exists := cmd.ConfigManager.GoWorkspaceExists(); !exists {
			cmd.SystemUtils.Logger.Verbose("go workspace does not exist yet, running 'go work init'")
			err := cmd.SystemUtils.Exec.GoCommand(cmd.Config.Runtime.ROOT, "work", "init")
			if err != nil {
				return err
			}
		} else {
			cmd.SystemUtils.Logger.Verbose("go workspace already exists, no need to create one")
		}
	} else if rootConfig.Strategy == "rewrite" {
		return errors.New("rewrite strategy unsupported yet")
	} else {
		return errors.New("invalid strategy '" + rootConfig.Strategy + "'")
	}

	if err := cmd.ConfigManager.WriteRootConfig(rootConfig); err != nil {
		return err
	} else {
		cmd.SystemUtils.Logger.Verbose("created monorepo configuration 'work.toml' at root")
	}

	// todo: check existence of modules folder (go.mod) to sanitize everything (create module.toml and make sure they are in the workspace)

	cmd.SystemUtils.Logger.Success("monorepo initialized at " + cmd.Config.Runtime.ROOT)

	return nil
}

func (cmd *Commands) List(c *cli.Context) error {
	if exists := cmd.ConfigManager.RootConfigExists(); !exists {
		return errors.New("monorepo not found at " + cmd.Config.Runtime.ROOT)
	}
	cmd.SystemUtils.Logger.Verbose("listing all modules")
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

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (cmd *Commands) Run(c *cli.Context) error {
	if exists := cmd.ConfigManager.RootConfigExists(); !exists {
		return errors.New("monorepo not found at " + cmd.Config.Runtime.ROOT)
	}

	scriptName := c.Args().Get(0)

	if scriptName == "" {
		return errors.New("no script name provided, usage: gorepo run [script_name]")
	} else {
		cmd.SystemUtils.Logger.Verbose("running script '" + scriptName + "'")
	}

	allowMissing := c.Bool("allow-missing")
	cmd.SystemUtils.Logger.Verbose("value for flag allowMissing: " + strconv.FormatBool(allowMissing))
	dryRun := c.Bool("dry-run")
	cmd.SystemUtils.Logger.Verbose("value for flag dryRun:       " + strconv.FormatBool(dryRun))
	targets := strings.Split(c.String("target"), ",")
	cmd.SystemUtils.Logger.Verbose("value for flag target:       " + strings.Join(targets, ","))
	for _, target := range targets {
		if target == "root" && len(targets) > 1 {
			return errors.New("cannot run script in root and in modules at the same time, you're being too fancy")
		}
	}

	if targets[0] == "root" {
		cmd.SystemUtils.Logger.Verbose("running script in root not supported yet")
	} else {
		allModules, err := cmd.ConfigManager.GetModules()
		if err != nil {
			return err
		}

		var modules []struct {
			RelativePath string
			ModuleConfig ModuleConfig
		}

		for _, module := range allModules {
			if targets[0] == "all" || contains(targets, module.ModuleConfig.Name) {
				modules = append(modules, module)
			}
		}

		// check all modules have the script
		cmd.SystemUtils.Logger.Verbose("checking if all modules have the script")
		var modulesWithoutScript []string
		for _, module := range modules {
			if _, ok := module.ModuleConfig.Scripts[scriptName]; !ok || module.ModuleConfig.Scripts[scriptName] == "" {
				modulesWithoutScript = append(modulesWithoutScript, module.ModuleConfig.Name)
			}
		}
		if len(modulesWithoutScript) == len(modules) {
			return errors.New("not running script, because it is missing in all modules")
		} else if len(modulesWithoutScript) > 0 && !allowMissing {
			return errors.New("not running script, because it is missing in following modules '" + scriptName + "' :" + strings.Join(modulesWithoutScript, ", "))
		} else if len(modulesWithoutScript) > 0 && allowMissing {
			cmd.SystemUtils.Logger.Verbose("script is missing in following modules (but flag allowMissing was passed) '" + scriptName + "' :" + strings.Join(modulesWithoutScript, ", "))
		} else {
			cmd.SystemUtils.Logger.Verbose("all modules have the script")
		}

		// execute them
		for _, module := range modules {
			path := filepath.Join(cmd.Config.Runtime.ROOT, module.RelativePath)
			script := module.ModuleConfig.Scripts[scriptName]
			if script == "" {
				cmd.SystemUtils.Logger.Info("script is empty, skipping")
				continue
			}
			cmd.SystemUtils.Logger.Info("running script " + scriptName + " in module " + module.ModuleConfig.Name)
			if dryRun {
				continue
			}
			if err := cmd.SystemUtils.Exec.BashCommand(path, script); err != nil {
				return err
			}
		}
	}

	return nil
}

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
