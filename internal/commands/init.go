package commands

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"gorepo-cli/internal/config"
	"os"
	"path/filepath"
	"strings"
)

func (cmd *Commands) Init(c *cli.Context) error {
	if exists := cmd.ConfigManager.RootConfigExists(); exists {
		return errors.New("monorepo already exists at " + cmd.Config.Runtime.ROOT)
	}

	rootConfig := config.RootConfig{
		Name:     c.Args().Get(0),
		Version:  "0.1.0",
		Strategy: "workspace",
		Vendor:   true,
	}

	// ask name
	if rootConfig.Name == "" {
		reader := bufio.NewReader(os.Stdin)
		defaultName := filepath.Base(cmd.Config.Runtime.ROOT)
		cmd.SystemUtils.Logger.Default("Enter the monorepo name: " + cmd.SystemUtils.Logger.ApplyInfoColor(fmt.Sprintf("default: ("+defaultName+")")) + " ")
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
	cmd.SystemUtils.Logger.Info("Using go workspace strategy by default (no other option for now)")

	// ask if should vendor
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Do you want to vendor dependencies? (y/n) " + cmd.SystemUtils.Logger.ApplyInfoColor("default: (y)") + " ")
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
