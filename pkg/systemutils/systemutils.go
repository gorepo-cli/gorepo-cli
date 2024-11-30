package systemutils

import (
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
	"os/exec"
)

// todo: add interface

type SystemUtils struct {
	Fs     Fs
	Exec   Exec
	Logger *LevelLogger
}

func NewSystemUtils() SystemUtils {
	return SystemUtils{
		Fs:     Fs{},
		Exec:   Exec{},
		Logger: NewLevelLogger(),
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
	cmd := exec.Command("/bin/sh", "-c", script)
	cmd.Dir = absolutePath // Set the working directory
	cmd.Stdout = os.Stdout // Redirect standard output to the parent process
	cmd.Stderr = os.Stderr // Redirect standard error to the parent process
	// Run the command and handle errors
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run command in %s: %w", absolutePath, err)
	}
	return nil
}

// LevelLogger is a logger that logs messages of different levels
type LevelLogger struct {
	*log.Logger
}

var (
	FatalColor   = color.New(color.FgRed).SprintFunc()
	WarningColor = color.New(color.FgYellow).SprintFunc()
	VerboseColor = color.New(color.FgHiBlack).SprintFunc()
	SuccessColor = color.New(color.FgGreen).SprintFunc()
	InfoColor    = color.New(color.FgCyan).SprintFunc()
)

func NewLevelLogger() *LevelLogger {
	return &LevelLogger{Logger: log.New(os.Stdout, "", 0)}
}

func (l *LevelLogger) FatalLn(msg string) {
	l.Println(FatalColor(msg))
}

func (l *LevelLogger) WarningLn(msg string) {
	l.Logger.Println(WarningColor(msg))
}

func (l *LevelLogger) VerboseLn(msg string) {
	l.Logger.Println(VerboseColor(msg))
}

func (l *LevelLogger) SuccessLn(msg string) {
	l.Logger.Println(SuccessColor(msg))
}

func (l *LevelLogger) InfoLn(msg string) {
	l.Logger.Println(InfoColor(msg))
}

func (l *LevelLogger) DefaultLn(msg string) {
	l.Logger.Println(msg)
}

func (l *LevelLogger) Default(msg string) {
	_, _ = l.Writer().Write([]byte(msg))
}
