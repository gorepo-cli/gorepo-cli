package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Runtime Runtime
	Static  Static
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

func NewConfig(su SystemUtils) (cfg Config, err error) {
	cfg.Static = Static{
		MaxRecursion:   7,
		RootFileName:   "work.toml",
		ModuleFileName: "module.toml",
	}
	cfg.Runtime = Runtime{}
	if wd, err := os.Getwd(); err == nil {
		cfg.Runtime.WD = wd
	} else {
		return cfg, err
	}
	if root, err := GetRoot(cfg, su); err == nil {
		cfg.Runtime.ROOT = root
	} else {
		return cfg, err
	}
	return cfg, nil
}

func GetRoot(cfg Config, su SystemUtils) (root string, err error) {
	currentDir := cfg.Runtime.WD
	if currentDir == "" {
		return "", fmt.Errorf("empty_wd")
	}
	for i := 0; i <= cfg.Static.MaxRecursion; i++ {
		filePath := filepath.Join(currentDir, cfg.Static.RootFileName)
		if exists := su.Fs.FileExists(filePath); exists {
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
