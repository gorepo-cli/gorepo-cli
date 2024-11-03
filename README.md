<p align="center">
    <picture>
        <img style="margin-bottom:0;" width="130" src="./gorepo.png" alt="logo">
    </picture>
    <h1 align="center">GOREPO-CLI</h1>
</p>

<p align="center">
    A CLI tool to manage golang monorepos.
</p>

<p align="center">
    /!\ Not nearly a v1
</p>

## Philosophy

- The CLI should be super dumb to use and transparent regarding what it does - --verbose should log everything
- The CLI would favorite flat structure but should be able to handle nested folders
- Commands should be ran from anywhere within the monorepo (avoiding the use of cd or having many terminals)
- No need to edit the go.work file, go.mod or gorepo.toml files manually
- Deleting the monorepo files (gorepo.toml files and .gorepo) should be enough to remove the CLI
- The CLI should be non intrusive and should not modify anything unless approved explicitely by the user

## Brainstorm commands

- `gorepo init` to create a new monorepo
    - -> ask for the type of monorepo (workspace vs module rewrite)
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