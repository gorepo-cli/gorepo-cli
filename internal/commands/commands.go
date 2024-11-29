package commands

import (
	"gorepo-cli/internal/config"
	"gorepo-cli/pkg/systemutils"
)

type Commands struct {
	SystemUtils   systemutils.SystemUtils
	Config        config.Config
	ConfigManager config.ConfigManager
}

func NewCommands(su systemutils.SystemUtils, cfg config.Config, manager config.ConfigManager) *Commands {
	return &Commands{
		SystemUtils:   su,
		Config:        cfg,
		ConfigManager: manager,
	}
}
