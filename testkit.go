package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func NewMockConfig(su SystemUtils, root, wd string) (cfg Config, err error) {
	cfg.Static = StaticConfig{
		MaxRecursion:   7,
		RootFileName:   "work.toml",
		ModuleFileName: "module.toml",
	}
	cfg.Runtime = RuntimeConfig{
		WD:   wd,
		ROOT: root,
	}
	cfg.su = su
	return cfg, nil
}

type MockFs struct {
	Files map[string][]byte
}

func NewMockFs() *MockFs {
	return &MockFs{
		Files: map[string][]byte{},
	}
}

func (m *MockFs) Exists(path string) bool {
	_, exists := m.Files[path]
	return exists
}

func (m *MockFs) Read(path string) ([]byte, error) {
	if data, exists := m.Files[path]; exists {
		return data, nil
	}
	return nil, os.ErrNotExist
}

func (m *MockFs) Write(path string, content []byte) error {
	m.Files[path] = content
	return nil
}

func (m *MockFs) Walk(root string, walkFn filepath.WalkFunc) error {
	//for path := range m.Files {
	//	info := mockFileInfo{
	//		name:    filepath.Base(path),
	//		size:    int64(len(m.Files[path])),
	//		mode:    0644,
	//		modTime: mockTime{},
	//		isDir:   false,
	//		sys:     nil,
	//	}
	//	if err := walkFn(path, info, nil); err != nil {
	//		return err
	//	}
	//}
	return nil
}

type mockFileInfo struct {
	name    string
	size    int64
	mode    fs.FileMode
	modTime mockTime
	isDir   bool
	sys     interface{}
}

func (m mockFileInfo) Name() string      { return m.name }
func (m mockFileInfo) Size() int64       { return m.size }
func (m mockFileInfo) Mode() fs.FileMode { return m.mode }
func (m mockFileInfo) ModTime() mockTime { return m.modTime }
func (m mockFileInfo) IsDir() bool       { return m.isDir }
func (m mockFileInfo) Sys() interface{}  { return m.sys }

type mockTime struct{}

func (mockTime) Unix() int64            { return 0 }
func (mockTime) String() string         { return "mockTime" }
func (mockTime) IsZero() bool           { return true }
func (mockTime) Before(t mockTime) bool { return false }

func (m *MockFs) Output(path string, perm os.FileMode) map[string][]byte {
	return m.Files
}

/////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////

type MockExec struct {
	Commands []MockCommand
}

func NewMockExec() *MockExec {
	return &MockExec{
		Commands: []MockCommand{},
	}
}

type MockCommand struct {
	Dir     string
	Command string
	Output  string
	Err     error
}

func (m *MockExec) GoCommand(dir string, args ...string) error {
	cmd := strings.Join(args, " ")
	m.Commands = append(m.Commands, MockCommand{
		Dir:     dir,
		Command: "go " + cmd,
	})
	return nil
}

func (m *MockExec) BashCommand(absolutePath, script string) error {
	m.Commands = append(m.Commands, MockCommand{
		Dir:     absolutePath,
		Command: script,
	})
	return nil
}

func (m *MockExec) Output() []MockCommand {
	return m.Commands
}

/////////////////////////////////////////////////////////////////

type MockLogger struct {
	Messages []string
}

func NewMockLogger() MockLogger {
	return MockLogger{
		Messages: []string{},
	}
}

func (l *MockLogger) FatalLn(msg string) {
	l.Messages = append(l.Messages, "FATAL: "+msg)
}

func (l *MockLogger) WarningLn(msg string) {
	l.Messages = append(l.Messages, "WARNING: "+msg)
}

func (l *MockLogger) VerboseLn(msg string) {
	l.Messages = append(l.Messages, "VERBOSE: "+msg)
}

func (l *MockLogger) SuccessLn(msg string) {
	l.Messages = append(l.Messages, "SUCCESS: "+msg)
}

func (l *MockLogger) InfoLn(msg string) {
	l.Messages = append(l.Messages, "INFO: "+msg)
}

func (l *MockLogger) DefaultLn(msg string) {
	l.Messages = append(l.Messages, "DEFAULT: "+msg)
}

func (l *MockLogger) Default(msg string) {
	l.Messages = append(l.Messages, "DEFAULT: "+msg)
}

func (l *MockLogger) Output() []string {
	return l.Messages
}
