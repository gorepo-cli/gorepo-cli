package main

import (
	"gorepo-cli/internal/cli"
	"gorepo-cli/pkg/systemutils"
	"os"
)

func main() {
	su := systemutils.NewSystemUtils()
	if err := cli.Exec(); err != nil {
		su.Logger.Fatal(err.Error())
		os.Exit(1)
	}
}
