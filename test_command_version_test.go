package main

import (
	"github.com/urfave/cli/v2"
	"testing"
)

func TestCommandVersion(t *testing.T) {
	t.Run("should return the version", func(t *testing.T) {
		mockLogger := NewMockLogger()
		su := NewSystemUtils(NewMockFs(), &MockExec{}, &mockLogger)
		cfg, err := NewMockConfig(su, "/root", "/root")
		if err != nil {
			t.Fatal(err)
		}
		if err != nil {
			t.Fatal(err)
		}
		cmd := NewCommands(su, cfg)
		err = cmd.Version(&cli.Context{})
		if err != nil {
			t.Fatal(err)
		}
		logs := mockLogger.Output()
		if logs[0] != "DEFAULT: dev" {
			t.Fatalf("Expected %s, got %s", "dev", logs[0])
		}
	})
}
