package config

import (
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"gorepo-cli/pkg/systemutils"
	"os"
	"path/filepath"
	"sort"
)

type Config struct {
	Runtime Runtime
	Static  Static
	su      systemutils.SystemUtils
}

type Runtime struct {
	WD   string // Working directory, folder where cli was executed
	ROOT string // Root of the monorepo
}

type Static struct {
	MaxRecursion   int    // Max recursion depth to search for monorepo root
	RootFileName   string // File name to identify the monorepo
	ModuleFileName string // File name to identify a module
}

/////////////////////////

func NewConfig(su systemutils.SystemUtils) (cfg Config, err error) {
	cfg.Static = Static{
		MaxRecursion:   7,
		RootFileName:   "work.toml",
		ModuleFileName: "module.toml",
	}
	cfg.Runtime = Runtime{}
	cfg.su = su
	if wd, err := os.Getwd(); err == nil {
		cfg.Runtime.WD = wd
	} else {
		return cfg, err
	}
	if root, err := getRootPath(cfg); err == nil {
		cfg.Runtime.ROOT = root
	} else {
		return cfg, err
	}
	return cfg, nil
}

func getRootPath(cfg Config) (root string, err error) {
	currentDir := cfg.Runtime.WD
	if currentDir == "" {
		return "", fmt.Errorf("empty_wd")
	}
	for i := 0; i <= cfg.Static.MaxRecursion; i++ {
		filePath := filepath.Join(currentDir, cfg.Static.RootFileName)
		if exists := cfg.su.Fs.FileExists(filePath); exists {
			return currentDir, nil
		}
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			return cfg.Runtime.WD, nil
		}
		currentDir = parentDir
	}
	return "", fmt.Errorf("root_not_found")
}

/////////////////////////
// Toml files
/////////////////////////

type RootConfig struct {
	Name     string            `toml:"name"`
	Version  string            `toml:"version"`
	Strategy string            `toml:"strategy"` // workspace / rewrites (unsupported yet)
	Vendor   bool              `toml:"vendor"`   // vendor or not
	Scripts  map[string]string `toml:"scripts"`
}

type ModuleConfig struct {
	Name         string            `toml:"-"` // name of the folder, added at runtime
	RelativePath string            `toml:"-"` // relative path to the root, added at runtime
	Scripts      map[string]string `toml:"scripts"`
}

func (c *Config) RootConfigExists() bool {
	filePath := filepath.Join(c.Runtime.ROOT, c.Static.RootFileName)
	return c.su.Fs.FileExists(filePath)
}

func (c *Config) LoadRootConfig() (cfg RootConfig, err error) {
	file, err := c.su.Fs.ReadFile(filepath.Join(c.Runtime.ROOT, c.Static.RootFileName))
	if err != nil {
		return cfg, err
	}
	err = toml.Unmarshal(file, &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}

func (c *Config) WriteRootConfig(rootConfig RootConfig) (err error) {
	configStr, err := toml.Marshal(rootConfig)
	if err != nil {
		return err
	}
	filePath := filepath.Join(c.Runtime.ROOT, c.Static.RootFileName)
	return c.su.Fs.WriteFile(filePath, configStr)
}

func (c *Config) GoWorkspaceExists() bool {
	filePath := filepath.Join(c.Runtime.ROOT, "go.work")
	return c.su.Fs.FileExists(filePath)
}

func (c *Config) GetModules() (modules []ModuleConfig, err error) {
	currentPath := c.Runtime.ROOT

	err = filepath.Walk(currentPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			exists := c.su.Fs.FileExists(filepath.Join(path, c.Static.ModuleFileName))
			if exists {
				relativePath, err := filepath.Rel(c.Runtime.ROOT, path)
				if err != nil {
					return err
				}
				moduleConfig, err := c.LoadModuleConfig(relativePath)
				if err != nil {
					return err
				}
				modules = append(modules, moduleConfig)
			}
		}

		return nil
	})

	if err != nil {
		c.su.Logger.Warning(err.Error())
		return modules, err
	}

	sort.Slice(modules, func(i, j int) bool {
		return modules[i].Name < modules[j].Name
	})

	return modules, nil
}

func (c *Config) LoadModuleConfig(relativePath string) (cfg ModuleConfig, err error) {
	path := filepath.Join(c.Runtime.ROOT, relativePath, c.Static.ModuleFileName)
	file, err := c.su.Fs.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	err = toml.Unmarshal(file, &cfg)
	if err != nil {
		return cfg, err
	}
	cfg.Name = filepath.Base(relativePath)
	cfg.RelativePath = relativePath
	return cfg, nil
}

func (c *Config) WriteModuleConfig(modConfig ModuleConfig, path, name string) (err error) {
	//configStr, err := toml.Marshal(modConfig)
	//if err != nil {
	//	return err
	//}
	//// todo: modules can be nested
	//filePath := filepath.Join(c.Config.Runtime.ROOT, c.Config.Static.ModuleFileName)
	return nil //c.SystemUtils.Fs.WriteFile(filePath, configStr)
}
