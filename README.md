<p align="center">
    <picture>
        <img style="margin-bottom:0;" width="130" src="./doc/assets/gorepo.png" alt="logo">
    </picture>
    <h1 align="center">GOREPO-CLI</h1>
</p>

<p align="center">
    A CLI to manage go monorepos.
</p>

<p align="center">
    /!\ Not nearly a v1
</p>

## Philosophy

- The CLI should be dumb to use. The dumber the better.
- The CLI should allow running all CI/CD commands for all modules at the root (test, lint, build, etc.)
- The CLI should favorite a flat monorepo structure, but should be able to handle nested folders
- The CLI should be non-intrusive and should not modify anything unless approved explicitely by the user
- The CLI should allow the user to add modules from templates, and add custom templates
- The CLI could be controlling docker, git and git hooks for the user
- The CLI could also handle incremental builds and deployments, given the user configures a storage for artifacts

## Pre-requisites

To use the CLI, you must have go installed since it runs go commands.

## Getting started as a dev

- Run `make build` to create bin/gorepo
- Add the bin folder to your PATH
- As a result, you can now run `gorepo` from anywhere

### Example on Fedora:
```bash
vim ~/.bashrc

# add this:
export PATH="$PATH:/home/my_name/Repositories/gorepo-cli/bin"

# refresh the terminal
source ~/.bashrc
```

## Documentation

This documentation is maintained as we go. It only contains information that 
is relevant to the actual version.

### To create a new monorepo at the current working directory

```bash
# Navigate to the root of your monorepo and run:
> gorepo init
```

This will create a work.toml file at the root, and a go.work file if one does not exist.
Currently, go workspaces are the only strategy available to build a monorepo with the CLI.

Note you can not nest monorepos.

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