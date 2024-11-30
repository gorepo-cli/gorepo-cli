package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/pelletier/go-toml/v2"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// SystemUtils contains side effect utilities that interact with the system
type SystemUtils struct {
	Fs     FsI
	Exec   ExecI
	Logger Llog
}

func NewSystemUtils(fs FsI, x ExecI, l Llog) SystemUtils {
	return SystemUtils{
		Fs:     fs,
		Exec:   x,
		Logger: l,
	}
}

// FsI defines methods to interact with the filesystem
type FsI interface {
	Exists(path string) bool
	Write(path string, content []byte) error
	Read(path string) ([]byte, error)
}

// Fs implements FsI
type Fs struct{}

var _ FsI = &Fs{}

func (fs *Fs) Exists(path string) (exists bool) {
	_, err := os.Stat(path)
	return err == nil
}

func (fs *Fs) Write(path string, content []byte) (err error) {
	return os.WriteFile(path, content, 0644)
}

func (fs *Fs) Read(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// ExecI defines methods to run commands
type ExecI interface {
	GoCommand(dir string, args ...string) error
	BashCommand(absolutePath, script string) error
}

// Exec implements ExecI
type Exec struct{}

var _ ExecI = &Exec{}

// GoCommand runs a go command in a given directory
func (x *Exec) GoCommand(absolutePath string, args ...string) (err error) {
	cmd := exec.Command("go", args...)
	cmd.Dir = absolutePath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run command: %w\nOutput: %s", err, string(output))
	}
	return nil
}

// BashCommand runs a bash script in a given directory
func (x *Exec) BashCommand(absolutePath, script string) (err error) {
	if _, err := os.Stat(absolutePath); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", absolutePath)
	}
	cmd := exec.Command("/bin/sh", "-c", script)
	cmd.Dir = absolutePath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run command in %s: %w", absolutePath, err)
	}
	return nil
}

// Llog is a logger that can log different levels
type Llog struct {
	*log.Logger
}

var (
	FatalColor   = color.New(color.FgRed).SprintFunc()
	WarningColor = color.New(color.FgYellow).SprintFunc()
	VerboseColor = color.New(color.FgHiBlack).SprintFunc()
	SuccessColor = color.New(color.FgGreen).SprintFunc()
	InfoColor    = color.New(color.FgCyan).SprintFunc()
)

func NewLevelLogger() *Llog {
	return &Llog{Logger: log.New(os.Stdout, "", 0)}
}

func (l *Llog) FatalLn(msg string) {
	l.Println(FatalColor(msg))
}

func (l *Llog) WarningLn(msg string) {
	l.Logger.Println(WarningColor(msg))
}

func (l *Llog) VerboseLn(msg string) {
	l.Logger.Println(VerboseColor(msg))
}

func (l *Llog) SuccessLn(msg string) {
	l.Logger.Println(SuccessColor(msg))
}

func (l *Llog) InfoLn(msg string) {
	l.Logger.Println(InfoColor(msg))
}

func (l *Llog) DefaultLn(msg string) {
	l.Logger.Println(msg)
}

func (l *Llog) Default(msg string) {
	_, _ = l.Writer().Write([]byte(msg))
}

// Config contains and manages configuration for the monorepo
type Config struct {
	Runtime RuntimeConfig
	Static  StaticConfig
	su      SystemUtils
}

type RuntimeConfig struct {
	WD   string // Working directory, folder where cli was executed
	ROOT string // Root of the monorepo
}

type StaticConfig struct {
	MaxRecursion   int    // Max recursion depth to search for monorepo root
	RootFileName   string // File name to identify the monorepo
	ModuleFileName string // File name to identify a module
}

func NewConfig(su SystemUtils) (cfg Config, err error) {
	cfg.Static = StaticConfig{
		MaxRecursion:   7,
		RootFileName:   "work.toml",
		ModuleFileName: "module.toml",
	}
	cfg.Runtime = RuntimeConfig{}
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
		return "", fmt.Errorf("no working directory")
	}
	for i := 0; i <= cfg.Static.MaxRecursion; i++ {
		filePath := filepath.Join(currentDir, cfg.Static.RootFileName)
		if exists := cfg.su.Fs.Exists(filePath); exists {
			return currentDir, nil
		}
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			return cfg.Runtime.WD, nil
		}
		currentDir = parentDir
	}
	return "", fmt.Errorf("root not found")
}

// RootConfig contains the configuration of the monorepo
type RootConfig struct {
	Name     string            `toml:"name"`
	Version  string            `toml:"version"`
	Strategy string            `toml:"strategy"` // workspace / rewrites (unsupported yet)
	Vendor   bool              `toml:"vendor"`   // vendor or not
	Scripts  map[string]string `toml:"scripts"`
}

// ModuleConfig contains the configuration of a module
type ModuleConfig struct {
	Name         string            `toml:"-"` // name of the folder, added at runtime
	RelativePath string            `toml:"-"` // relative path to the root, added at runtime
	Scripts      map[string]string `toml:"scripts"`
}

func (c *Config) RootConfigExists() bool {
	filePath := filepath.Join(c.Runtime.ROOT, c.Static.RootFileName)
	return c.su.Fs.Exists(filePath)
}

func (c *Config) LoadRootConfig() (cfg RootConfig, err error) {
	file, err := c.su.Fs.Read(filepath.Join(c.Runtime.ROOT, c.Static.RootFileName))
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
	return c.su.Fs.Write(filePath, configStr)
}

func (c *Config) GoWorkspaceExists() bool {
	filePath := filepath.Join(c.Runtime.ROOT, "go.work")
	return c.su.Fs.Exists(filePath)
}

func (c *Config) GetModules() (modules []ModuleConfig, err error) {
	currentPath := c.Runtime.ROOT
	err = filepath.Walk(currentPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			exists := c.su.Fs.Exists(filepath.Join(path, c.Static.ModuleFileName))
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
		c.su.Logger.WarningLn(err.Error())
		return modules, err
	}
	sort.Slice(modules, func(i, j int) bool {
		return modules[i].Name < modules[j].Name
	})
	return modules, nil
}

func (c *Config) LoadModuleConfig(relativePath string) (cfg ModuleConfig, err error) {
	path := filepath.Join(c.Runtime.ROOT, relativePath, c.Static.ModuleFileName)
	file, err := c.su.Fs.Read(path)
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

func (c *Config) WriteModuleConfig(modConfig ModuleConfig, relativePath, name string) (err error) {
	// todo
	return nil
}

// Commands contains the CLI commands
type Commands struct {
	SystemUtils SystemUtils
	Config      Config
}

func NewCommands(su SystemUtils, cfg Config) *Commands {
	return &Commands{
		SystemUtils: su,
		Config:      cfg,
	}
}

func (cmd *Commands) Init(c *cli.Context) error {
	if exists := cmd.Config.RootConfigExists(); exists {
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
		cmd.SystemUtils.Logger.Default("Enter the monorepo name: " + InfoColor(fmt.Sprintf("default: ("+defaultName+")")) + " ")
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
	cmd.SystemUtils.Logger.InfoLn("Using go workspace strategy by default (no other option for now)")

	// ask if should vendor
	reader := bufio.NewReader(os.Stdin)
	cmd.SystemUtils.Logger.Default("Do you want to vendor dependencies? (y/n) " + InfoColor("default: (y)") + " ")
	response, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}
	rootConfig.Vendor = strings.TrimSpace(response) == "y"

	// handle go workspace
	if rootConfig.Strategy == "workspace" {
		if exists := cmd.Config.GoWorkspaceExists(); !exists {
			cmd.SystemUtils.Logger.VerboseLn("go workspace does not exist yet, running 'go work init'")
			err := cmd.SystemUtils.Exec.GoCommand(cmd.Config.Runtime.ROOT, "work", "init")
			if err != nil {
				return err
			}
		} else {
			cmd.SystemUtils.Logger.VerboseLn("go workspace already exists, no need to create one")
		}
	} else if rootConfig.Strategy == "rewrite" {
		return errors.New("rewrite strategy unsupported yet")
	} else {
		return errors.New("invalid strategy '" + rootConfig.Strategy + "'")
	}

	if err := cmd.Config.WriteRootConfig(rootConfig); err != nil {
		return err
	} else {
		cmd.SystemUtils.Logger.VerboseLn("created monorepo configuration 'work.toml' at root")
	}

	// todo: check existence of modules folder (go.mod) to sanitize everything (create module.toml and make sure they are in the workspace)

	cmd.SystemUtils.Logger.SuccessLn("monorepo initialized at " + cmd.Config.Runtime.ROOT)

	return nil
}

func (cmd *Commands) List(c *cli.Context) error {
	if exists := cmd.Config.RootConfigExists(); !exists {
		return errors.New("monorepo not found at " + cmd.Config.Runtime.ROOT)
	}
	modules, err := cmd.Config.GetModules()
	if err != nil {
		return err
	}
	if len(modules) == 0 {
		cmd.SystemUtils.Logger.InfoLn("no modules found")
	} else {
		for _, module := range modules {
			cmd.SystemUtils.Logger.DefaultLn(module.Name)
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
	if exists := cmd.Config.RootConfigExists(); !exists {
		return errors.New("monorepo not found at " + cmd.Config.Runtime.ROOT)
	}

	scriptName := c.Args().Get(0)

	if scriptName == "" {
		return errors.New("no script name provided, usage: gorepo run [script_name]")
	} else {
		cmd.SystemUtils.Logger.VerboseLn("running script '" + scriptName + "'")
	}

	allowMissing := c.Bool("allow-missing")
	cmd.SystemUtils.Logger.VerboseLn("value for flag allowMissing: " + strconv.FormatBool(allowMissing))
	dryRun := c.Bool("dry-run")
	cmd.SystemUtils.Logger.VerboseLn("value for flag dryRun:       " + strconv.FormatBool(dryRun))
	targets := strings.Split(c.String("target"), ",")
	cmd.SystemUtils.Logger.VerboseLn("value for flag target:       " + strings.Join(targets, ","))
	for _, target := range targets {
		if target == "root" && len(targets) > 1 {
			return errors.New("cannot run script in root and in modules at the same time, you're being too fancy")
		}
	}

	if targets[0] == "root" {
		cmd.SystemUtils.Logger.VerboseLn("running script in root not supported yet")
	} else {
		allModules, err := cmd.Config.GetModules()
		if err != nil {
			return err
		}

		var modules []ModuleConfig

		for _, module := range allModules {
			if targets[0] == "all" || contains(targets, module.Name) {
				modules = append(modules, module)
			}
		}

		// check all modules have the script
		cmd.SystemUtils.Logger.VerboseLn("checking if all modules have the script")
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
			cmd.SystemUtils.Logger.VerboseLn("script is missing in following modules (but flag allowMissing was passed) '" + scriptName + "' :" + strings.Join(modulesWithoutScript, ", "))
		} else {
			cmd.SystemUtils.Logger.VerboseLn("all modules have the script")
		}

		// execute them
		for _, module := range modules {
			path := filepath.Join(cmd.Config.Runtime.ROOT, module.RelativePath)
			script := module.Scripts[scriptName]
			if script == "" {
				cmd.SystemUtils.Logger.InfoLn("script is empty, skipping")
				continue
			}
			cmd.SystemUtils.Logger.InfoLn("running script " + scriptName + " in module " + module.Name)
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

var version = "dev"

func (cmd *Commands) Version(c *cli.Context) error {
	cmd.SystemUtils.Logger.DefaultLn(version)
	return nil
}

func (cmd *Commands) Debug(c *cli.Context) error {
	cmd.SystemUtils.Logger.InfoLn("===================")
	cmd.SystemUtils.Logger.InfoLn("RUNTIME_CONFIG")
	cmd.SystemUtils.Logger.InfoLn("===================")
	cmd.SystemUtils.Logger.DefaultLn("WD_(COMMAND_RAN_FROM)........" + cmd.Config.Runtime.WD)
	cmd.SystemUtils.Logger.DefaultLn("ROOT (OF THE MONOREPO)......." + cmd.Config.Runtime.ROOT)
	cmd.SystemUtils.Logger.DefaultLn("MONOREPO EXISTS (AT ROOT)...." +
		strconv.FormatBool(cmd.Config.RootConfigExists()))

	cmd.SystemUtils.Logger.InfoLn("===================")
	cmd.SystemUtils.Logger.InfoLn("STATIC_CONFIG")
	cmd.SystemUtils.Logger.InfoLn("===================")
	cmd.SystemUtils.Logger.DefaultLn("MAX RECURSION................" + strconv.Itoa(cmd.Config.Static.MaxRecursion))
	cmd.SystemUtils.Logger.DefaultLn("ROOT FILE NAME..............." + cmd.Config.Static.RootFileName)
	cmd.SystemUtils.Logger.DefaultLn("MODULE FILE NAME............." + cmd.Config.Static.ModuleFileName)

	if cmd.Config.RootConfigExists() {
		cmd.SystemUtils.Logger.InfoLn("===================")
		cmd.SystemUtils.Logger.InfoLn("ROOT_CONFIG")
		cmd.SystemUtils.Logger.InfoLn("===================")

		cfg, err := cmd.Config.LoadRootConfig()
		if err != nil {
			return err
		}

		cmd.SystemUtils.Logger.DefaultLn("NAME.........." + cfg.Name)
		cmd.SystemUtils.Logger.DefaultLn("VERSION......." + cfg.Version)
		cmd.SystemUtils.Logger.DefaultLn("STRATEGY......" + cfg.Strategy)
		cmd.SystemUtils.Logger.DefaultLn("VENDOR........" + strconv.FormatBool(cfg.Vendor))

		modules, err := cmd.Config.GetModules()
		if err != nil {
			return err
		}

		cmd.SystemUtils.Logger.DefaultLn("N_MODULES....." + strconv.Itoa(len(modules)))

		if len(modules) > 0 {
			cmd.SystemUtils.Logger.InfoLn("===================")
			cmd.SystemUtils.Logger.InfoLn("MODULES_CONFIG")
			cmd.SystemUtils.Logger.InfoLn("===================")
		}

		for _, module := range modules {
			cmd.SystemUtils.Logger.InfoLn("MODULE " + module.Name)
			cmd.SystemUtils.Logger.DefaultLn("MODULE_NAME........ " + module.Name)
			cmd.SystemUtils.Logger.DefaultLn("MODULE_PATH........ " + module.RelativePath)
			if len(module.Scripts) > 0 {
				cmd.SystemUtils.Logger.DefaultLn("COMMANDS........")
				for k, v := range module.Scripts {
					cmd.SystemUtils.Logger.DefaultLn("  " + k + " -> " + v)
				}
			}
		}
	}
	return nil
}

// Run runs the CLI application
func Run() (err error) {
	su := NewSystemUtils(&Fs{}, &Exec{}, *NewLevelLogger())
	cfg, err := NewConfig(su)
	if err != nil {
		return err
	}
	cmd := NewCommands(su, cfg)
	app := &cli.App{
		Name:  "GOREPO",
		Usage: "A CLI tool to manage Go monorepos",
		Commands: []*cli.Command{
			{
				Name:   "init",
				Usage:  "Initialize a new monorepo at the working directory",
				Action: cmd.Init,
			},
			{
				Name:   "list",
				Usage:  "List all modules in the monorepo",
				Action: cmd.List,
			},
			{
				Name:   "run",
				Usage:  "Run a command in a given scope (all modules, some modules, at root)",
				Action: cmd.Run,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "target",
						Value: "all",
						Usage: "Target root or specific modules (comma separated)",
					},
					&cli.BoolFlag{
						Name:  "dry-run",
						Value: false,
						Usage: "Print the commands that would be executed",
					},
					&cli.BoolFlag{
						Name:  "allow-missing",
						Value: false,
						Usage: "Run the scripts in the modules that have it, even if it is missing in some",
					},
					//&cli.BoolFlag{
					//	Name:  "parallel",
					//	Value: false,
					//	Usage: "Run the scripts in parallel",
					//},
				},
			},
			{
				Name:   "version",
				Usage:  "Print the version of the monorepo",
				Action: cmd.Version,
			},
			{
				Name:   "debug",
				Usage:  "Gives information about the configuration",
				Action: cmd.Debug,
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "verbose",
				Usage: "Enable verbose logging for all commands",
				Value: false,
			},
		},
	}
	return app.Run(os.Args)
}

func main() {
	su := NewSystemUtils(&Fs{}, &Exec{}, *NewLevelLogger())
	if err := Run(); err != nil {
		su.Logger.FatalLn(err.Error())
		os.Exit(1)
	}
}