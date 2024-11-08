package main

import (
	"github.com/pelletier/go-toml/v2"
	"os"
	"path/filepath"
)

type RootConfig struct {
	Name     string
	Version  string
	Strategy string // workspace / rewrites (unsupported yet)
	Vendor   bool   // vendor or not
}

type ModuleConfig struct {
	Name string
}

type ConfigManager struct {
	SystemUtils SystemUtils
	Config      Config
}

func NewConfigManager(su SystemUtils, cfg Config) ConfigManager {
	return ConfigManager{
		SystemUtils: su,
		Config:      cfg,
	}
}

func (c *ConfigManager) RootConfigExists() bool {
	filePath := filepath.Join(c.Config.Runtime.ROOT, c.Config.Static.RootFileName)
	if _, err := os.Stat(filePath); err == nil {
		return true
	}
	return false
}

func (c *ConfigManager) LoadRootConfig() (cfg RootConfig, err error) {
	return cfg, nil
}

func (c *ConfigManager) WriteRootConfig(rootConfig RootConfig) (err error) {
	configStr, err := toml.Marshal(rootConfig)
	if err != nil {
		return err
	}
	filePath := filepath.Join(c.Config.Runtime.ROOT, c.Config.Static.RootFileName)
	return c.SystemUtils.Fs.WriteFile(filePath, configStr)
}

func (c *ConfigManager) GoWorkspaceExists() bool {
	filePath := filepath.Join(c.Config.Runtime.ROOT, "go.work")
	if _, err := os.Stat(filePath); err == nil {
		return true
	}
	return false
}

func (c *ConfigManager) LoadModuleConfig(mod string) (cfg ModuleConfig, err error) {
	return cfg, nil
}

func (c *ConfigManager) WriteModuleConfig(modConfig string) (err error) {
	//configStr, err := toml.Marshal(modConfig)
	//if err != nil {
	//	return err
	//}
	//// todo: modules can be nested
	//filePath := filepath.Join(c.Config.Runtime.ROOT, c.Config.Static.ModuleFileName)
	return nil //c.SystemUtils.Fs.WriteFile(filePath, configStr)
}
