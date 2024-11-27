
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