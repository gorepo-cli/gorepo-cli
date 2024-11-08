package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	fatal   = color.New(color.FgRed).SprintFunc()     // Red for errors
	warning = color.New(color.FgYellow).SprintFunc()  // Yellow for warnings
	verbose = color.New(color.FgHiBlack).SprintFunc() // Gray for general messages
	comment = color.New(color.FgCyan).SprintFunc()    // Cyan for comments
	success = color.New(color.FgGreen).SprintFunc()   // Green for success messages
	info    = color.New(color.FgBlue).SprintFunc()    // Blue for info messages
)

func main() {
	env, err := NewEnv()
	if err != nil {
		log.Fatal(fatal(err))
	}
	commands := NewCommands(env)
	app := &cli.App{
		Name:  "GOREPO-CLI",
		Usage: "A CLI tool to manage Go monorepos",
		Commands: []*cli.Command{
			{
				Name:   "init",
				Usage:  "Initialize a new monorepo at the working directory",
				Action: commands.Init,
			},
			//{
			//	Name:   "add",
			//	Usage:  "Add a new module to the monorepo",
			//	Action: addModule,
			//	Flags: []cli.Flag{
			//		&cli.StringFlag{
			//			Name:  "template",
			//			Usage: "Template to use for creating the module",
			//		},
			//	},
			//},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(fatal(err))
	}
}

///////////////////////////////////////////////////////
// Commands
///////////////////////////////////////////////////////

type Commands struct {
	Env Env
}

func NewCommands(env Env) *Commands {
	return &Commands{
		Env: env,
	}
}

func (commands *Commands) Init(c *cli.Context) error {
	fmt.Println(verbose("Initializing a new monorepo..."))

	// Check monorepo config existence
	if exists := DoesWorkspaceExist(commands.Env); exists {
		errMsg := "Monorepo already exists at " + commands.Env.ROOT
		return errors.New(errMsg)
	}

	config := WorkConfig{
		Name:     c.Args().Get(0),
		Version:  "0.1.0",
		Strategy: "workspace",
		Vendor:   true,
	}

	// ask name
	if config.Name == "" {
		reader := bufio.NewReader(os.Stdin)
		base := filepath.Base(commands.Env.ROOT)
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
		if exists := DoesGoWorkspaceExist(commands.Env); !exists {
			RunGoCommand(commands.Env.ROOT, "work", "init")
		}
	} else if config.Strategy == "rewrite" {
		return errors.New("Unsupported yet")
	} else {
		return errors.New("Invalid strategy")
	}

	// Create monorepo config

	// check existence of modules folder (go.mod) to sanitize everything (create module.toml and make sure they are in the workspace)

	fmt.Println("ended init")
	return nil
}

///////////////////////////////////////////////////////
// Env
///////////////////////////////////////////////////////

type Env struct {
	WD               string // Working directory, folder where cli was executed
	ROOT             string // Root of the monorepo, if exists
	MaxRecursion     int    // Max recursion depth to search for monorepo root
	MonorepoFileName string // File name to identify the monorepo
	ModuleFileName   string // File name to identify a module
}

func NewEnv() (env Env, err error) {
	env = Env{
		MaxRecursion:     7,
		MonorepoFileName: "work.toml",
		ModuleFileName:   "module.toml",
	}
	if wd, err := os.Getwd(); err == nil {
		env.WD = wd
	} else {
		return env, err
	}
	if root, err := GetRoot(env); err == nil {
		env.ROOT = root
	} else {
		return env, err
	}
	return env, nil
}

// Returns the root of the monorepo
func GetRoot(env Env) (root string, err error) {
	currentDir := env.WD
	if currentDir == "" {
		return "", fmt.Errorf("empty_dir")
	}
	for i := 0; i <= env.MaxRecursion; i++ {
		filePath := filepath.Join(currentDir, env.MonorepoFileName)
		if _, err := os.Stat(filePath); err == nil {
			return currentDir, nil
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			return env.WD, nil
		}
		currentDir = parentDir
	}
	return "", fmt.Errorf("not_found")
}

// configuration from the work.toml file
type WorkConfig struct {
	Name     string
	Version  string
	Strategy string // workspace / rewrites (unsupported yet)
	Vendor   bool   // vendor or not
}

func DoesWorkspaceExist(env Env) bool {
	filePath := filepath.Join(env.ROOT, env.MonorepoFileName)
	if _, err := os.Stat(filePath); err == nil {
		return true
	}
	return false
}

func DoesGoWorkspaceExist(env Env) bool {
	filePath := filepath.Join(env.ROOT, "go.work")
	if _, err := os.Stat(filePath); err == nil {
		return true
	}
	return false
}

// configuration from the module.toml file
type ModuleConfig struct {
	Name string
}

func RunGoCommand(dir string, args ...string) error {
	cmd := exec.Command("go", args...)

	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run command: %w\nOutput: %s", err, string(output))
	}

	fmt.Printf("Command output: %s\n", string(output))
	return nil
}
