package main

import (
	"flag"
	"github.com/pelletier/go-toml/v2"
	"github.com/urfave/cli/v2"
	"testing"
)

func TestCommandInit(t *testing.T) {
	t.Run("should return an error if a work.toml already exists at root", func(t *testing.T) {
		rootConfigBytes, _ := toml.Marshal(RootConfig{
			Name: "my-monorepo",
		})
		tk, _ := NewTestKit("/some/path/root", map[string][]byte{
			"/some/path/root/work.toml": rootConfigBytes,
		}, nil, nil)
		mockContext := cli.NewContext(&cli.App{
			Name:  "test-app",
			Usage: "This is just a test",
		}, flag.NewFlagSet("test", flag.ContinueOnError), nil)
		err := tk.cmd.Init(mockContext)
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
		if err.Error() != "monorepo already exists at /some/path/root" {
			t.Fatalf("expected 'work.toml already exists at root', got %s", err.Error())
		}
	})
	t.Run("should create a new work.toml if there is no such file at root", func(t *testing.T) {
		tk, _ := NewTestKit("/some/path/root", map[string][]byte{}, map[string]bool{
			"Do you want to vendor dependencies?": true,
		}, map[string]string{
			"What is the monorepo name?": "",
		})
		mockContext := cli.NewContext(&cli.App{
			Name:  "test-app",
			Usage: "This is just a test",
		}, flag.NewFlagSet("test", flag.ContinueOnError), nil)
		err := tk.cmd.Init(mockContext)
		if err != nil {
			t.Fatal("expected an error, got", err)
		}
		files := tk.MockFs.Output()
		if files["/some/path/root/work.toml"] == nil {
			t.Fatal("expected a non-nil value, got nil")
		}
	})
	t.Run("should create a go.work file if it is missing", func(t *testing.T) {
		// todo
	})
	t.Run("should use the default name if the user does not select one", func(t *testing.T) {
		// todo
	})
	t.Run("should use the name passed as an argument", func(t *testing.T) {
		// todo
	})
	t.Run("should use the name provided by the user if they answer the question about it", func(t *testing.T) {
		// todo
	})
	t.Run("should set vendoring to true if the user selected it", func(t *testing.T) {
		// todo
	})
}
