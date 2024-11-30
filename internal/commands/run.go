package commands

import (
	"errors"
	"github.com/urfave/cli/v2"
	"gorepo-cli/internal/config"
	"path/filepath"
	"strconv"
	"strings"
)

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (cmd *Commands) Run(c *cli.Context) error {
	if exists := cmd.Config.RootConfigExists(); !exists {
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
		allModules, err := cmd.Config.GetModules()
		if err != nil {
			return err
		}

		var modules []config.ModuleConfig

		for _, module := range allModules {
			if targets[0] == "all" || contains(targets, module.Name) {
				modules = append(modules, module)
			}
		}

		// check all modules have the script
		cmd.SystemUtils.Logger.Verbose("checking if all modules have the script")
		var modulesWithoutScript []string
		for _, module := range modules {
			if _, ok := module.Scripts[scriptName]; !ok || module.Scripts[scriptName] == "" {
				modulesWithoutScript = append(modulesWithoutScript, module.Name)
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
			script := module.Scripts[scriptName]
			if script == "" {
				cmd.SystemUtils.Logger.Info("script is empty, skipping")
				continue
			}
			cmd.SystemUtils.Logger.Info("running script " + scriptName + " in module " + module.Name)
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
