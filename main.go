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

// SystemUtils contains utilities to interact with the system
type SystemUtils struct {
	Fs     FsI
	Exec   ExecI
	Logger LlogI
	Os     OsI
}

// NewSystemUtils returns an instance of SystemUtils
func NewSystemUtils(fs FsI, x ExecI, l LlogI, o OsI) *SystemUtils {
	return &SystemUtils{
		Fs:     fs,
		Exec:   x,
		Logger: l,
		Os:     o,
	}
}

// FsI defines methods to interact with the filesystem
type FsI interface {
	Exists(path string) bool
	Read(path string) ([]byte, error)
	Write(path string, content []byte) error
	Walk(root string, walkFn filepath.WalkFunc) error
}

// Fs implements FsI
type Fs struct{}

var _ FsI = &Fs{}

// Exists checks if a file exists
func (fs *Fs) Exists(path string) (exists bool) {
	_, err := os.Stat(path)
	return err == nil
}

// Read reads and return a file
func (fs *Fs) Read(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// Write writes content to a file
func (fs *Fs) Write(path string, content []byte) (err error) {
	return os.WriteFile(path, content, 0644)
}

// Walk walks the filesystem
func (fs *Fs) Walk(root string, walkFn filepath.WalkFunc) (err error) {
	return filepath.Walk(root, walkFn)
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

// LlogI defines methods to log messages
type LlogI interface {
	FatalLn(msg string)
	WarningLn(msg string)
	VerboseLn(msg string)
	SuccessLn(msg string)
	InfoLn(msg string)
	DefaultLn(msg string)
	Default(msg string)
}

// Llog implements LlogI
type Llog struct {
	*log.Logger
}

var _ LlogI = &Llog{}

// NewLevelLogger returns an instance of Llog
func NewLevelLogger() *Llog {
	return &Llog{Logger: log.New(os.Stdout, "", 0)}
}

var (
	FatalColor   = color.New(color.FgRed).SprintFunc()
	WarningColor = color.New(color.FgYellow).SprintFunc()
	VerboseColor = color.New(color.FgHiBlack).SprintFunc()
	SuccessColor = color.New(color.FgGreen).SprintFunc()
	InfoColor    = color.New(color.FgCyan).SprintFunc()
)

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

// OsI defines methods to interact with the operating system
type OsI interface {
	GetWd() (dir string, err error)
	AskBool(question, choices, defaultValue string, logger LlogI) (response bool, err error)
	AskString(question, choices, defaultValue string, logger LlogI) (response string, err error)
}

// Os implements OsI
type Os struct{}

var _ OsI = &Os{}

// GetWd returns the working directory
func (o *Os) GetWd() (dir string, err error) {
	return os.Getwd()
}

// AskBool asks a question and returns a boolean
func (o *Os) AskBool(question, choices, defaultValue string, logger LlogI) (response bool, err error) {
	questionFormated := question
	choicesFormated := ""
	if choices != "" {
		choicesFormated = InfoColor("(" + choices + ")")
	}
	defaultValueFormated := ""
	if defaultValue != "" {
		defaultValueFormated = VerboseColor("default: " + defaultValue)
	}
	logger.Default(questionFormated + " " + choicesFormated + " " + defaultValueFormated + ": ")
	reader := bufio.NewReader(os.Stdin)
	responseStr, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	responseStr = strings.TrimSpace(strings.ToLower(responseStr))
	if responseStr == "" {
		responseStr = defaultValue
	}
	return responseStr == "y" || responseStr == "yes", nil
}

// AskString asks a question and returns a string
func (o *Os) AskString(question, choices, defaultValue string, logger LlogI) (response string, err error) {
	questionFormated := question
	choicesFormated := ""
	if choices != "" {
		choicesFormated = InfoColor("(" + choices + ")")
	}
	defaultValueFormated := ""
	if defaultValue != "" {
		defaultValueFormated = VerboseColor("default: " + defaultValue)
	}
	logger.Default(questionFormated + " " + choicesFormated + " " + defaultValueFormated + ": ")
	reader := bufio.NewReader(os.Stdin)
	responseStr, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	responseStr = strings.TrimSpace(responseStr)
	if responseStr == "" {
		responseStr = defaultValue
	}
	return responseStr, nil
}

// Config contains configuration of the monorepo
type Config struct {
	Static  StaticConfig
	Runtime RuntimeConfig
	su      *SystemUtils
}

// StaticConfig contains static configuration of the monorepo
type StaticConfig struct {
	MaxRecursion   int    // Max recursion depth to search for monorepo root
	RootFileName   string // File name to identify the monorepo
	ModuleFileName string // File name to identify a module
}

// RuntimeConfig contains runtime variables
type RuntimeConfig struct {
	WD   string // Working directory, folder where cli was executed
	ROOT string // Root of the monorepo
}

// RootManipulation defines methods to manipulate the root configuration
type RootManipulation interface {
	RootConfigExists() bool
	LoadRootConfig() (cfg RootConfig, err error)
	WriteRootConfig(rootConfig RootConfig) (err error)
}

var _ RootManipulation = &Config{}

// ModuleManipulation defines methods to manipulate the module configuration
type ModuleManipulation interface {
	GetModules(targets, exclude []string) (modules []ModuleConfig, err error)
	LoadModuleConfig(relativePath string) (cfg ModuleConfig, err error)
	WriteModuleConfig(modConfig ModuleConfig, relativePath, name string) (err error)
}

var _ ModuleManipulation = &Config{}

// ConfigHelpers defines methods to help with the configuration
type ConfigHelpers interface {
	GoWorkspaceExists() bool
}

var _ ConfigHelpers = &Config{}

// NewConfig returns an instance of Config
func NewConfig(su *SystemUtils) (cfg *Config, err error) {
	cfg = &Config{}
	cfg.Static = StaticConfig{
		MaxRecursion:   7,
		RootFileName:   "work.toml",
		ModuleFileName: "module.toml",
	}
	cfg.Runtime = RuntimeConfig{}
	cfg.su = su
	if wd, err := su.Os.GetWd(); err == nil {
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

func getRootPath(cfg *Config) (root string, err error) {
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
	Strategy string            `toml:"strategy"` // workspace / rewrites (unsupported)
	Vendor   bool              `toml:"vendor"`   // vendor or not (unsupported)
	Scripts  map[string]string `toml:"scripts"`
}

// ModuleConfig contains the configuration of a module
type ModuleConfig struct {
	Name         string            `toml:"-"` // name of the folder, added at runtime
	RelativePath string            `toml:"-"` // relative path to the root, added at runtime
	Scripts      map[string]string `toml:"scripts"`
}

// RootConfigExists checks if a file work.toml exists at the root
func (c *Config) RootConfigExists() bool {
	filePath := filepath.Join(c.Runtime.ROOT, c.Static.RootFileName)
	return c.su.Fs.Exists(filePath)
}

// LoadRootConfig loads the root configuration of the monorepo
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

// WriteRootConfig writes the root configuration of the monorepo
func (c *Config) WriteRootConfig(rootConfig RootConfig) (err error) {
	configStr, err := toml.Marshal(rootConfig)
	if err != nil {
		return err
	}
	filePath := filepath.Join(c.Runtime.ROOT, c.Static.RootFileName)
	return c.su.Fs.Write(filePath, configStr)
}

// GoWorkspaceExists checks if a file go.work exists at the root
func (c *Config) GoWorkspaceExists() bool {
	filePath := filepath.Join(c.Runtime.ROOT, "go.work")
	return c.su.Fs.Exists(filePath)
}

// GetModules returns all modules in the monorepo in alphabetical order
func (c *Config) GetModules(targets, exclude []string) (modules []ModuleConfig, err error) {
	// validation
	for _, target := range targets {
		if target == "root" && len(targets) > 1 {
			return nil, errors.New("cannot run script in root and in modules at the same time, you're being too greedy, run the command twice")
		} else if target == "all" && len(targets) > 1 {
			return nil, errors.New("cannot run script in all modules and in specific modules, non sense")
		}
	}
	for _, excluded := range exclude {
		if excluded == "all" {
			return nil, errors.New("excluding all modules makes no sense")
		} else if excluded == "root" {
			return nil, errors.New("excluding root is the default behaviour, no need to specify it")
		}
	}
	// walk
	currentPath := c.Runtime.ROOT
	err = c.su.Fs.Walk(currentPath, func(path string, info os.FileInfo, err error) error {
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
				if (targets[0] == "all" || contains(targets, moduleConfig.Name)) && !contains(exclude, moduleConfig.Name) {
					modules = append(modules, moduleConfig)
				}
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

// LoadModuleConfig loads the configuration of a module
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

// WriteModuleConfig writes the configuration of a module
func (c *Config) WriteModuleConfig(modConfig ModuleConfig, relativePath, name string) (err error) {
	// todo
	return nil
}

// Commands contains the actual CLI commands
type Commands struct {
	SystemUtils *SystemUtils
	Config      *Config
}

// NewCommands returns an instance of Commands
func NewCommands(su *SystemUtils, cfg *Config) *Commands {
	return &Commands{
		SystemUtils: su,
		Config:      cfg,
	}
}

////////////////////////////////////////
// CLI COMMANDS
////////////////////////////////////////

// Init implements `gorepo init`
func (cmd *Commands) Init(c *cli.Context) error {
	if exists := cmd.Config.RootConfigExists(); exists {
		return errors.New("monorepo already exists at " + cmd.Config.Runtime.ROOT)
	}

	verbose := c.Bool("verbose")

	rootConfig := RootConfig{
		Name:     c.Args().Get(0),
		Version:  "0.1.0",
		Strategy: "workspace",
		Vendor:   true,
	}

	// ask name
	if rootConfig.Name == "" {
		defaultName := filepath.Base(cmd.Config.Runtime.ROOT)
		nameResponse, err := cmd.SystemUtils.Os.AskString("What is the monorepo name?", "", defaultName, cmd.SystemUtils.Logger)
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		rootConfig.Name = nameResponse
	}

	// ask strategy
	cmd.SystemUtils.Logger.InfoLn("Using go workspace strategy by default (no other option for now)")

	// ask if should vendor
	if vendorResponse, err := cmd.SystemUtils.Os.AskBool(
		"Do you want to vendor dependencies?", "y/n", "y", cmd.SystemUtils.Logger); err == nil {
		rootConfig.Vendor = vendorResponse
	} else {
		return fmt.Errorf("failed to read input: %w", err)
	}

	// handle go workspace
	if rootConfig.Strategy == "workspace" {
		if exists := cmd.Config.GoWorkspaceExists(); !exists {
			if verbose {
				cmd.SystemUtils.Logger.VerboseLn("go workspace does not exist yet, running 'go work init'")
			}
			err := cmd.SystemUtils.Exec.GoCommand(cmd.Config.Runtime.ROOT, "work", "init")
			if err != nil {
				return err
			}
		} else {
			if verbose {
				cmd.SystemUtils.Logger.VerboseLn("go workspace already exists, no need to create one")
			}
			// todo: handle vendoring
		}
	} else if rootConfig.Strategy == "rewrite" {
		return errors.New("rewrite strategy unsupported yet")
	} else {
		return errors.New("invalid strategy '" + rootConfig.Strategy + "'")
	}

	if err := cmd.Config.WriteRootConfig(rootConfig); err != nil {
		return err
	} else {
		if verbose {
			cmd.SystemUtils.Logger.VerboseLn("created monorepo configuration 'work.toml' at root")
		}
	}

	// todo: check existence of modules folder (go.mod) to sanitize everything (create module.toml and make sure they are in the workspace)

	cmd.SystemUtils.Logger.SuccessLn("monorepo successfully initialized at " + cmd.Config.Runtime.ROOT)

	return nil
}

// List implements `gorepo list`
func (cmd *Commands) List(c *cli.Context) error {
	if exists := cmd.Config.RootConfigExists(); !exists {
		return errors.New("monorepo not found at " + cmd.Config.Runtime.ROOT)
	}
	modules, err := cmd.Config.GetModules([]string{"all"}, []string{})
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

// Execute implements `gorepo execute`
func (cmd *Commands) Execute(c *cli.Context) error {
	if exists := cmd.Config.RootConfigExists(); !exists {
		return errors.New("monorepo not found at " + cmd.Config.Runtime.ROOT)
	}

	verbose := c.Bool("verbose")
	if verbose {
		cmd.SystemUtils.Logger.VerboseLn("verbose mode enabled")
	}

	scriptName := c.Args().Get(0)
	if scriptName == "" {
		return errors.New("no script name provided, usage: gorepo run [script_name]")
	} else {
		if verbose {
			cmd.SystemUtils.Logger.VerboseLn("running script '" + scriptName + "'")
		}
	}

	allowMissing := c.Bool("allow-missing")
	if verbose {
		cmd.SystemUtils.Logger.VerboseLn("value for flag allowMissing: " + strconv.FormatBool(allowMissing))
	}

	targets := strings.Split(c.String("target"), ",")
	if verbose {
		cmd.SystemUtils.Logger.VerboseLn("value for flag target:       " + strings.Join(targets, ","))
	}

	exclude := strings.Split(c.String("exclude"), ",")
	if verbose {
		cmd.SystemUtils.Logger.VerboseLn("value for flag exclude:      " + strings.Join(exclude, ","))
	}

	// logic

	if targets[0] == "root" {
		cmd.SystemUtils.Logger.WarningLn("running script in root not supported yet")
		// implement here and return
		return nil
	}

	modules, err := cmd.Config.GetModules(targets, exclude)
	if err != nil {
		return err
	}

	if len(modules) == 0 {
		return errors.New("no modules found")
	}

	// check all modules have the script
	if verbose && !allowMissing {
		cmd.SystemUtils.Logger.VerboseLn("checking if all modules have the script")
	}
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
		if verbose {
			cmd.SystemUtils.Logger.VerboseLn("script is missing in following modules (but flag allowMissing was passed) '" + scriptName + "' :" + strings.Join(modulesWithoutScript, ", "))
		}
	} else {
		if verbose {
			cmd.SystemUtils.Logger.VerboseLn("all modules have the script")
		}
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
		if err := cmd.SystemUtils.Exec.BashCommand(path, script); err != nil {
			return err
		}
	}

	return nil
}

// FmtCI implements `gorepo fmt-ci`
func (cmd *Commands) FmtCI(c *cli.Context) error {
	if exists := cmd.Config.RootConfigExists(); !exists {
		return errors.New("monorepo not found at " + cmd.Config.Runtime.ROOT)
	}

	verbose := c.Bool("verbose")
	if verbose {
		cmd.SystemUtils.Logger.VerboseLn("verbose mode enabled")
	}

	targets := strings.Split(c.String("target"), ",")
	if verbose {
		cmd.SystemUtils.Logger.VerboseLn("value for flag target:       " + strings.Join(targets, ","))
	}

	exclude := strings.Split(c.String("exclude"), ",")
	if verbose {
		cmd.SystemUtils.Logger.VerboseLn("value for flag exclude:      " + strings.Join(exclude, ","))
	}

	if targets[0] == "root" {
		return errors.New("running fmt in root is not supported")
	}

	modules, err := cmd.Config.GetModules(targets, exclude)
	if err != nil {
		return err
	}

	script := "if [ -n \"$(gofmt -l .)\" ]; then exit 1; fi"

	for _, module := range modules {
		path := filepath.Join(cmd.Config.Runtime.ROOT, module.RelativePath)
		if err := cmd.SystemUtils.Exec.BashCommand(path, script); err != nil {
			return errors.New("fmt failed in module " + module.Name)
		}
	}

	return nil
}

// VetCI implements `gorepo vet-ci`
func (cmd *Commands) VetCI(c *cli.Context) error {
	experimental := c.Bool("experimental")
	if !experimental {
		return errors.New("this is an experimental feature, use --experimental flag to enable it")
	}

	if exists := cmd.Config.RootConfigExists(); !exists {
		return errors.New("monorepo not found at " + cmd.Config.Runtime.ROOT)
	}

	verbose := c.Bool("verbose")
	if verbose {
		cmd.SystemUtils.Logger.VerboseLn("verbose mode enabled")
	}

	targets := strings.Split(c.String("target"), ",")
	if verbose {
		cmd.SystemUtils.Logger.VerboseLn("value for flag target:       " + strings.Join(targets, ","))
	}

	exclude := strings.Split(c.String("exclude"), ",")
	if verbose {
		cmd.SystemUtils.Logger.VerboseLn("value for flag exclude:      " + strings.Join(exclude, ","))
	}

	if targets[0] == "root" {
		return errors.New("running fmt in root is not supported")
	}

	modules, err := cmd.Config.GetModules(targets, exclude)
	if err != nil {
		return err
	}

	script := "go vet . 2>&1 | grep -q . && exit 1"

	for _, module := range modules {
		path := filepath.Join(cmd.Config.Runtime.ROOT, module.RelativePath)
		if err := cmd.SystemUtils.Exec.BashCommand(path, script); err != nil {
			return errors.New("vet failed in module " + module.Name)
		}
	}

	return nil
}

// version is injected at build time
var version = "dev"

// Version implements `gorepo version`
func (cmd *Commands) Version(c *cli.Context) error {
	cmd.SystemUtils.Logger.DefaultLn(version)
	return nil
}

// Diagnostic implements `gorepo debug`
func (cmd *Commands) Diagnostic(c *cli.Context) error {
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

		modules, err := cmd.Config.GetModules([]string{"all"}, []string{})
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

// Cli runs the CLI application
func Cli() (err error) {
	su := NewSystemUtils(&Fs{}, &Exec{}, NewLevelLogger(), &Os{})
	cfg, err := NewConfig(su)
	if err != nil {
		return err
	}
	cmd := NewCommands(su, cfg)
	executionFlags := []cli.Flag{
		&cli.StringFlag{
			Name:  "target",
			Value: "all",
			Usage: "Target specific modules (comma separated)",
		},
		&cli.StringFlag{
			Name:  "exclude",
			Value: "",
			Usage: "Exclude specific modules (comma separated)",
		},
		&cli.BoolFlag{
			Name:  "allow-missing",
			Value: false,
			Usage: "Allow executing the scripts, even if some module don't have it",
		},
	}
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
				Name:  "add",
				Usage: "NOT IMPLEMENTED - Add a new module to the monorepo",
			},
			{
				Name:   "list",
				Usage:  "List all modules of the monorepo",
				Action: cmd.List,
			},
			{
				Name:   "execute",
				Usage:  "Execute a script across targeted modules",
				Action: cmd.Execute,
				Flags:  executionFlags,
			},
			{
				Name:   "fmt-ci",
				Usage:  "Breaks if targeted modules are not formatted",
				Action: cmd.FmtCI,
				Flags:  executionFlags,
			},
			{
				Name:   "vet-ci",
				Usage:  "[experimental] Breaks if targeted modules have vet issues",
				Action: cmd.VetCI,
				Flags:  executionFlags,
			},
			{
				Name:   "version",
				Usage:  "Print the version of the CLI",
				Action: cmd.Version,
			},
			{
				Name:   "diagnostic",
				Usage:  "Gives information about the configuration (mostly for internal use)",
				Hidden: true,
				Action: cmd.Diagnostic,
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "verbose",
				Usage: "Enable verbose logging for all commands",
				Value: false,
			},
			&cli.BoolFlag{
				Name:  "experimental",
				Value: false,
				Usage: "Experiment new features",
			},
		},
	}
	return app.Run(os.Args)
}

// main is the entry point
func main() {
	su := NewSystemUtils(&Fs{}, &Exec{}, NewLevelLogger(), &Os{})
	if err := Cli(); err != nil {
		su.Logger.FatalLn(err.Error())
		os.Exit(1)
	}
}
