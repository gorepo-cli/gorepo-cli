package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
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
		return errors.New("Monorepo already exists at " + cmd.Config.Runtime.ROOT)
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
		fmt.Print("Enter the monorepo name: " + comment(fmt.Sprintf("default: ("+base+")")) + " ")
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
	fmt.Println(comment("Using go workspace strategy by default (no other option for now)"))

	// ask if should vendor
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Do you want to vendor dependencies? (y/n) " + comment("default: (y)") + " ")
	response, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}
	config.Vendor = strings.TrimSpace(response) == "y"

	// Check go workspace existence
	if config.Strategy == "workspace" {
		if exists := cmd.ConfigManager.GoWorkspaceExists(); !exists {
			err := cmd.SystemUtils.Exec.GoCommand(cmd.Config.Runtime.ROOT, "work", "init")
			if err != nil {
				return err
			}
		}
	} else if config.Strategy == "rewrite" {
		return errors.New("Unsupported yet")
	} else {
		return errors.New("Invalid strategy")
	}

	if err := cmd.ConfigManager.WriteRootConfig(config); err != nil {
		return err
	}

	// todo: check existence of modules folder (go.mod) to sanitize everything (create module.toml and make sure they are in the workspace)

	cmd.SystemUtils.Logger.Success("Monorepo initialized at " + cmd.Config.Runtime.ROOT)

	return nil
}
