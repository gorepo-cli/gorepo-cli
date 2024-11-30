package systemutils

import (
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
	"os/exec"
)

// SystemUtils contains side effect utilities that interact with the system
type SystemUtils struct {
	Fs     FsInteractions
	Exec   Executor
	Logger *LevelLogger
}

func NewSystemUtils() SystemUtils {
	return SystemUtils{
		Fs:     &Fs{},
		Exec:   &Exec{},
		Logger: NewLevelLogger(),
	}
}

// FsInteractions defines methods to interact with the filesystem
type FsInteractions interface {
	FileExists(path string) bool
	WriteFile(path string, content []byte) error
	ReadFile(path string) ([]byte, error)
}

// Fs implements the FsInteractions interface
type Fs struct{}

var _ FsInteractions = &Fs{}

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

// Executor defines methods to run bash commands
type Executor interface {
	GoCommand(dir string, args ...string) error
	BashCommand(absolutePath, script string) error
}

// Exec implements the Executor interface
type Exec struct{}

var _ Executor = &Exec{}

// GoCommand runs a go command in a given directory
func (x *Exec) GoCommand(absolutePath string, args ...string) (err error) {
	cmd := exec.Command("go", args...)
	cmd.Dir = absolutePath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run command: %w\nOutput: %s", err, string(output))
	}
	return nil
}

// BashCommand runs a bash script in a given directory
func (x *Exec) BashCommand(absolutePath, script string) (err error) {
	if _, err := os.Stat(absolutePath); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", absolutePath)
	}
	cmd := exec.Command("/bin/sh", "-c", script)
	cmd.Dir = absolutePath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run command in %s: %w", absolutePath, err)
	}
	return nil
}

// LevelLogger is a logger that can log at different levels
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
