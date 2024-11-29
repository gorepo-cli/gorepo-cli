<p align="center">
    <picture>
        <img style="margin-bottom:0;" width="130" src="./assets/gorepo.png" alt="logo">
    </picture>
    <h1 align="center">GOREPO-CLI</h1>
</p>

<p align="center">
    A CLI to manage Go monorepos.
</p>

## Philosophy

- The CLI should be dumb to use
- The CLI should allow running all commands from anywhere because having to cd is just annoying
- The CLI should allow running CI/CD commands (test, lint, build, etc.) for all modules at once
- The CLI should handle both flat and nested structures (but common, make it flat, as flat as the earth)
- The CLI should be transparent to the user

### Following statements are for future work:

- The CLI should allow the user to add modules from templates
- The users should be able to create templates
- The creator of a template should be allowed to define template-scripts
- The CLI could also handle incremental builds, given the user configures a storage

## Disclaimer
- This is not nearly a v1
- Commit before running the CLI to see exactly what you are doing with it
- I only test Linux for now, macOS is probably ok, Windows is probably not

## Pre-requisites

To use the CLI, you must have go installed since it runs go commands.

## Getting started

Gorepo is not yet available on any package manager. You will need to build it yourself:
- Clone/download the repository
- Run `make build` to create bin/gorepo
- Add the bin folder to your PATH
- As a result, you can now run `gorepo` from anywhere
- Change code, build, test from anywhere, repeat

### Example on Linux:
```
vim ~/.bashrc

# add this:
export PATH="$PATH:/home/my_name/Repositories/gorepo-cli/bin"

# refresh the terminal
source ~/.bashrc
```

## Concepts

todo

## Reference

The reference documentation contains information that 
is relevant to the actual commited version. Future development should be in BRAINSTORMING.md

### gorepo init

#### Description

Creates a monorepo at the current work directory

#### Usage

```
gorepo init
```

#### Behaviour

This command creates two primary files:
- `work.toml` at the work directory
- `go.work` file if the strategy is set as 'workspace' one does not exist yet. This runs `go work init` behind the hood

### gorepo list

#### Description

Lists all modules in the monorepo

#### Usage

```
gorepo list
```

#### Behaviour

This command lists all modules in the monorepo, formally a module is a folder with a `module.toml` file in it.

### gorepo run

#### Description

Runs a script from a specified context

#### Usage

```
gorepo run [--target] [--dry-run] [--allow-missing] [script]
```

#### Parameters

- `script`: the name of the script to run
- `--target` (optional): the name of the comma-separated module(s) to run the script in, examples: `--target=root`, `--target=mod1,mod2`
- `--dry-run` (optional): prints the command that would be run without actually running it
- `--allow-missing` (optional): allows the script to run even if some of the targets does not have the script

#### Behaviour

This command runs all the scripts (bash scripts) defined in `work.toml` and `module.toml` files that are targeted.
By default, it will not run if one of the targeted module is missing the script.
Note it will run all or nothing. If one fails, it will not revert the operations that already ran.
