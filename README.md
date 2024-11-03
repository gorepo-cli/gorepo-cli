<p align="center">
    <picture>
        <img style="margin-bottom:0;" width="130" src="./gorepo.png" alt="logo">
    </picture>
    <h1 align="center">GOREPO-CLI</h1>
</p>


<p style="text-align:center;">
    A CLI tool to manage golang monorepos.
</p>

<p style="text-align:center;">
    /!\ Not nearly a v1
</p>

## Brainstorm commands

- `gorepo init` to create a new monorepo
    - -> ask for the type of monorepo (workspace vs module rewrite)
    - -> ask vendor or not
    - -> generate go.work
    - -> generate gorepo.yml
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
- `gorepo remove mod` to remove a module from the monorepo
- `gorepo list` to list all modules in the monorepo
- `gorepo rename mod new_name` to rename a module
- `gorepo use xxx install xxx` to install a dependency in a module
- `gorepo use xxx special command from the template`
    - `gorepo use xxx add service`
    - `gorepo use xxx add endpoint`
- `gorepo lint`
- `gorepo fmt`
- `gorepo test`
- `gorepo build`
- `gorepo start`
- `gorepo dev`