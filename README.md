<p align="center">
    <picture>
        <img style="margin-bottom:0;" width="130" src="./doc/assets/gorepo.png" alt="logo">
    </picture>
    <h1 align="center">GOREPO-CLI</h1>
</p>

<p align="center">
    A CLI to manage Go monorepos.
</p>

## Philosophy

- The CLI should be dumb to use. The dumber the better.
- The CLI should allow running all CI/CD commands for all modules at once (test, lint, build, etc.)
- The CLI should allow all commands to be be ran from anywhere in the monorepo 
- The CLI should favorite a flat monorepo structure, but should be able to handle nested folders
- The CLI should be non-intrusive and should not modify anything unless approved explicitely by the user
- The CLI should allow the user to add modules from templates, and add custom templates
- The CLI could be controlling docker, git and git hooks for the user
- The CLI could also handle incremental builds and deployments, given the user configures a storage

## Disclaimer
- This is not nearly a v1
- Everything can change from a day to another until it is a v1
- This must be considered unstable
- Commit before running the CLI to see exactly what you are doing with it
- I test it only on Linux for now, it should be ok for Mac, but Windows is not considered yet

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

## Contributing

- Open an issue, discuss a change
- Fork the repository
- Create a branch
- Make your changes, Commit, Push
- Create a pull request
- Wait for a review
- Merge your PR