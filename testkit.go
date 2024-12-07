package main

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// TestKit mocks all the dependencies that interacts with the system for testing
type TestKit struct {
	MockLogger *MockLogger
	MockFs     *MockFs
	MockExec   *MockExec
	MockOs     *MockOs
	su         *SystemUtils
	cfg        *Config
	cmd        *Commands
}

// NewTestKit creates a new TestKit
// wd: working directory from where the command is executed
// files: map of files with their content (pass nil if not needed)
// qABool: map of questions and answers for boolean questions (pass nil if not needed)
// qaString: map of questions and answers for string questions (pass nil if not needed)
func NewTestKit(wd string, files map[string][]byte, qABool map[string]bool, qaString map[string]string) (tk *TestKit, err error) {
	mockFs := NewMockFs(files)
	mockExec := NewMockExec()
	mockLogger := NewMockLogger()
	mockOs := NewMockOs(wd, qABool, qaString)
	su := NewSystemUtils(mockFs, mockExec, mockLogger, mockOs)
	cfg, err := NewConfig(su)
	if err != nil {
		return &TestKit{}, err
	}
	return &TestKit{
		MockLogger: mockLogger,
		MockFs:     mockFs,
		MockExec:   mockExec,
		MockOs:     mockOs,
		su:         su,
		cfg:        cfg,
		cmd:        NewCommands(su, cfg),
	}, nil
}

type MockFs struct {
	Files map[string][]byte
}

func NewMockFs(files map[string][]byte) *MockFs {
	return &MockFs{
		Files: files,
	}
}

func (m MockFs) Exists(path string) bool {
	_, exists := m.Files[path]
	return exists
}

func (m MockFs) Read(path string) ([]byte, error) {
	if data, exists := m.Files[path]; exists {
		return data, nil
	}
	return nil, os.ErrNotExist
}

func (m MockFs) Write(path string, content []byte) error {
	m.Files[path] = content
	return nil
}

func (m MockFs) Walk(root string, walkFn filepath.WalkFunc) error {
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

func (m MockFs) Output(path string, perm os.FileMode) map[string][]byte {
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

func (m MockExec) GoCommand(dir string, args ...string) error {
	cmd := strings.Join(args, " ")
	m.Commands = append(m.Commands, MockCommand{
		Dir:     dir,
		Command: "go " + cmd,
	})
	return nil
}

func (m MockExec) BashCommand(absolutePath, script string) error {
	m.Commands = append(m.Commands, MockCommand{
		Dir:     absolutePath,
		Command: script,
	})
	return nil
}

func (m MockExec) Output() []MockCommand {
	return m.Commands
}

/////////////////////////////////////////////////////////////////

type MockLogger struct {
	Messages []string
}

func NewMockLogger() *MockLogger {
	return &MockLogger{
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

/////////////////////////////////////////////////////////////////

type MockOs struct {
	Wd                     string
	QuestionsAnswersBool   map[string]bool
	QuestionsAnswersString map[string]string
}

func NewMockOs(wd string, qABool map[string]bool, qAString map[string]string) *MockOs {
	return &MockOs{
		Wd:                     wd,
		QuestionsAnswersBool:   qABool,
		QuestionsAnswersString: qAString,
	}
}

func (m *MockOs) GetWd() (string, error) {
	return m.Wd, nil
}

func (m *MockOs) AskBool(question, choices, defaultValue string, logger LlogI) (response bool, err error) {
	// todo: implement
	return false, errors.New("not implemented")
}

func (m *MockOs) AskString(question, choices, defaultValue string, logger LlogI) (response string, err error) {
	// todo: implement
	return "", errors.New("not implemented")
}
