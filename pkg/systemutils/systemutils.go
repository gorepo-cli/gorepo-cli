package systemutils

import (
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
	"os/exec"
)

type SystemUtils struct {
	Fs     Fs
	Exec   Exec
	Logger Logger
}

func NewSystemUtils() SystemUtils {
	return SystemUtils{
		Fs:   Fs{},
		Exec: Exec{},
		Logger: Logger{
			Logger: log.New(os.Stdout, "", 0),
		},
	}
}

type Fs struct{}

func (fs *Fs) FileExists(path string) (exists bool) {
	_, err := os.Stat(path)
	return err == nil
}

func (fs *Fs) WriteFile(path string, content []byte) (err error) {
	return os.WriteFile(path, content, 0644)
}

func (fs *Fs) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

type Exec struct{}

func (x *Exec) GoCommand(dir string, args ...string) (err error) {
	cmd := exec.Command("go", args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run command: %w\nOutput: %s", err, string(output))
	}
	return nil
}

func (x *Exec) BashCommand(absolutePath, script string) (err error) {
	// Validate that the directory exists
	if _, err := os.Stat(absolutePath); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", absolutePath)
	}
	// Create the command to run the script
	cmd := exec.Command("/bin/bash", "-c", script)
	cmd.Dir = absolutePath // Set the working directory
	cmd.Stdout = os.Stdout // Redirect standard output to the parent process
	cmd.Stderr = os.Stderr // Redirect standard error to the parent process
	// Run the command and handle errors
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run command in %s: %w", absolutePath, err)
	}
	return nil
}

var (
	fatal   = color.New(color.FgRed).SprintFunc()
	warning = color.New(color.FgYellow).SprintFunc()
	verbose = color.New(color.FgHiBlack).SprintFunc()
	success = color.New(color.FgGreen).SprintFunc()
	info    = color.New(color.FgCyan).SprintFunc()
)

type Logger struct {
	Logger *log.Logger
}

func (su *Logger) Fatal(msg string) {
	su.Logger.Println(fatal(msg))
}

func (su *Logger) Warning(msg string) {
	su.Logger.Println(warning(msg))
}

func (su *Logger) Verbose(msg string) {
	su.Logger.Println(verbose(msg))
}

func (su *Logger) Success(msg string) {
	su.Logger.Println(success(msg))
}

func (su *Logger) Info(msg string) {
	su.Logger.Println(info(msg))
}

func (su *Logger) ApplyInfoColor(msg string) string {
	return info(msg)
}

func (su *Logger) Default(msg string) {
	su.Logger.Println(msg)
}
