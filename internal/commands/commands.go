package commands

import (
	"gorepo-cli/internal/config"
	"gorepo-cli/pkg/systemutils"
)

type Commands struct {
	SystemUtils systemutils.SystemUtils
	Config      config.Config
}

func NewCommands(su systemutils.SystemUtils, cfg config.Config) *Commands {
	return &Commands{
		SystemUtils: su,
		Config:      cfg,
	}
}
