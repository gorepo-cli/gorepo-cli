package main

import (
	"github.com/urfave/cli/v2"
	"testing"
)

func TestCommandVersion(t *testing.T) {
	t.Run("should return the version", func(t *testing.T) {
		tk, err := NewTestKit("/root", "/root", nil)
		if err != nil {
			t.Fatal(err)
		}
		err = tk.cmd.Version(&cli.Context{})
		if err != nil {
			t.Fatal(err)
		}
		logs := tk.MockLogger.Output()
		if logs[0] != "DEFAULT: dev" {
			t.Fatalf("Expected %s, got %s", "dev", logs[0])
		}
	})
}
