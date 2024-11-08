package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path/filepath"
)

func main() {
	env, err := NewEnv()
	if err != nil {
		log.Fatal(err)
		return
	}
	commands := NewCommands(env)
	app := &cli.App{
		Name:  "GOREPO-CLI",
		Usage: "A CLI tool to manage Go monorepos",
		Commands: []*cli.Command{
			{
				Name:   "init",
				Usage:  "Initialize a new monorepo",
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
		log.Fatal(err)
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
	fmt.Println("Initializing a new monorepo...", commands.Env)
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

// configuration from the module.toml file
type ModuleConfig struct {
	Name string
}
