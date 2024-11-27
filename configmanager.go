package main

import (
	"github.com/pelletier/go-toml/v2"
	"os"
	"path/filepath"
	"strings"
)

type RootConfig struct {
	Name     string
	Version  string
	Strategy string // workspace / rewrites (unsupported yet)
	Vendor   bool   // vendor or not
}

type ModuleConfig struct {
	Name string // the module name is the folder name with no spaces
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
	file, err := c.SystemUtils.Fs.ReadFile(filepath.Join(c.Config.Runtime.ROOT, c.Config.Static.RootFileName))
	if err != nil {
		return cfg, err
	}
	err = toml.Unmarshal(file, &cfg)
	if err != nil {
		return cfg, err
	}
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

// loop to check if folder exists
func (c *ConfigManager) ModuleFolderExists(name string) bool {
	currentPath := c.Config.Runtime.ROOT

	found := false

	err := filepath.Walk(currentPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && info.Name() == name {
			found = true
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		c.SystemUtils.Logger.Warning(err.Error())
		return false
	}

	return found
}

func (c *ConfigManager) GetModules() (modules []struct {
	RelativePath string
	ModuleConfig ModuleConfig
}, err error) {
	currentPath := c.Config.Runtime.ROOT

	err = filepath.Walk(currentPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			exists := c.SystemUtils.Fs.FileExists(filepath.Join(path, c.Config.Static.ModuleFileName))
			if exists {
				moduleConfig, err := c.LoadModuleConfig(path)
				if err != nil {
					return err
				}
				relativePath, err := filepath.Rel(c.Config.Runtime.ROOT, path)
				modules = append(modules, struct {
					RelativePath string
					ModuleConfig ModuleConfig
				}{RelativePath: relativePath, ModuleConfig: moduleConfig})
			}
		}

		return nil
	})

	if err != nil {
		c.SystemUtils.Logger.Warning(err.Error())
		return modules, err
	}

	return modules, nil
}

func (c *ConfigManager) LoadModuleConfig(path string) (cfg ModuleConfig, err error) {
	file, err := c.SystemUtils.Fs.ReadFile(filepath.Join(path, c.Config.Static.ModuleFileName))
	if err != nil {
		return cfg, err
	}
	err = toml.Unmarshal(file, &cfg)
	if err != nil {
		return cfg, err
	}
	folderName := filepath.Base(path)
	cfg.Name = strings.ReplaceAll(folderName, " ", "")
	return cfg, nil
}

func (c *ConfigManager) WriteModuleConfig(modConfig ModuleConfig, path, name string) (err error) {
	//configStr, err := toml.Marshal(modConfig)
	//if err != nil {
	//	return err
	//}
	//// todo: modules can be nested
	//filePath := filepath.Join(c.Config.Runtime.ROOT, c.Config.Static.ModuleFileName)
	return nil //c.SystemUtils.Fs.WriteFile(filePath, configStr)
}
