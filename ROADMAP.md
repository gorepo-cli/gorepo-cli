

# Roadmap & Brainstorm

## New Commands

- add:    to add a module
- remove: to remove a module
- health: to check the health of the modules (or check), with --fix
- fmt
- vet
- test
- build   (check how to set priority)
- run     (check how to know the path + priority)

## New flags
- [executionFlags] parallel: to run the commands in parallel
- [global]         dry-run:  to show what would be done 

### Following statements are for future work:

- The CLI should allow the user to add modules from templates
- The users should be able to create templates
- The creator of a template should be allowed to define template-scripts
- The CLI could also handle incremental builds, given the user configures a storage

## Brainstorm for future commands

- `gorepo init` to create a new monorepo
    - -> ask for the type of monorepo (workspace vs rewrite)
    - -> ask vendor or not
    - -> generate go.work
    - -> generate gorepo.toml
    - -> generate .gitignore
    - -> generate docker-compose.yml (if user wants)
    - -> how many servers
- `gorepo add new_mod` to add a new module to the monorepo
    - -> generate the module go.mod
    - -> adds the module to go.work
    - -> generate the module new_mod/gorepo.yml
    - `--template` to add a module from a template
        - @echo
        - @templ
        - @nginx
        - @kafka
        - @ghcicd
- `gorepo remove mod` to remove a module from the monorepo
- `gorepo list` to list all modules in the monorepo
- `gorepo rename mod new_name` to rename a module
- `gorepo use xxx install xxx` to install a dependency in a module
- `gorepo use xxx lint` to install a dependency in a module
- `gorepo use xxx run` to install a dependency in a module
- ...
- `gorepo use xxx special command from the template`
    - `gorepo use xxx add service`
    - `gorepo use xxx add endpoint`
- `gorepo lint`
- `gorepo fmt`
- `gorepo test`
- `gorepo build` // with/without docker - push or not
- `gorepo start` (call what was built) option `--watch` (runs dev, if docker), option `--no-docker` (runs dev, without docker)
- `gorepo check` flag `--fix` - `gorepo build`
- `gorepo tidy`
- `gorepo tree` to display the tree of the monorepo
- `gorepo version`
- `gorepo update` to update the CLI
- `gorepo help` to display the help
- `gorepo upgrade` to upgrade the packages to the latest version

example of docker flows

```bash
gorepo docker build --module api
gorepo docker push --module api --registry my-docker-registry
gorepo docker deploy --module api --env production
gorepo docker compose up --detach
gorepo docker compose logs --follow
```

--verbose should log everything

## Toml example

work.toml
```toml
[monorepo]
name = "MyGoMonorepo"
version = "1.0"
monorepo_strategy = "workspace"
vendor = true/false

[scripts]
```

module.toml
```toml
[module]
name = "api"

[template]
name = "@echo"
......

[commands]
     run = "go run cmd/service-1/main.go"
     lint = "golangci-lint run"
     test = "go test ./..."
```

## 1

- `grpo version`
- `grpo init` -> init work.tml, go.work, vendor, gitignore (no git, no docker, no modules)
- `grpo add name-or-path`
- `grpo list`

## 2

- `grpo lint`
- `grpo fmt`
- `grpo test`

## 5

- `grpo remove name-or-path`
- `grpo rename name-or-path new-name`


```
ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_.@!#$%^&()[]{}'+,;=~
```



//{
//	Name:   "add",
//	Usage:  "Add a new module to the monorepo",
//	Action: commands.Add,
//	Flags: []cli.Flag{
//		&cli.BoolFlag{
//			Name:  "verbose",
//			Usage: "Enable verbose output",
//		},
//		&cli.StringFlag{
//			Name:  "template",
//			Usage: "Choose a template (not implemented)",
//		},
//	},
//},
//{}, // sanitize / lint / health / check


Add Context Support for Cancellation

Issue: Long-running operations cannot be cancelled by the user.

Recommendation: Pass context.Context to functions to handle cancellation and timeouts.

go
Copy code
func (cmd *Commands) Run(c *cli.Context) error {
ctx := c.Context
// Pass ctx to functions and check for cancellation
}

Provide Execution Summaries

Issue: Users don't receive a summary of the executed commands.

Recommendation: Collect and display a summary at the end of script execution.

go
Copy code
var failedModules []string
// ... (during execution)
if err != nil {
failedModules = append(failedModules, module.Name)
}
// After execution
if len(failedModules) > 0 {
cmd.SystemUtils.Logger.Warning("Scripts failed in modules: " + strings.Join(failedModules, ", "))
} else {
cmd.SystemUtils.Logger.Success("All scripts ran successfully.")
}