package subcommand

import (
	"testing"

	"bytes"
	"errors"
	"os"
	"strings"

	"github.com/Oppodelldog/git-commit-hook/config"
	"github.com/Oppodelldog/git-commit-hook/hook"
	"github.com/Oppodelldog/git-commit-hook/testhelper"
	"github.com/stretchr/testify/assert"
)

func TestTestCommand_Test_HappyPath(t *testing.T) {
	originArgs := os.Args
	os.Args = []string{"programm name", "test", "-m", "test commit message", "-b", `"feature/PROJECT-123"`, "-p", "test project"}
	defer func() { os.Args = originArgs }()
	defer testhelper.CleanupTestEnvironment(t)
	testhelper.PreapreTestEnvironment(t)

	test := NewTestCommand()
	test.stdoutWriter = bytes.NewBufferString("")

	res := test.Test()

	expectedOutput := `
testing configuration '/tmp/git-commit-hook/git-commit-hook.yaml':
project        : test project
branch         : "feature/PROJECT-123"
commit message : test commit message

would generate the following commit message:
test commit message
`
	assert.Exactly(t, strings.TrimLeft(expectedOutput, "\n"), test.stdoutWriter.(*bytes.Buffer).String())
	assert.Exactly(t, 0, res)
}

func TestTestCommand_Test_TooFewParameters_PrintUsage(t *testing.T) {
	originArgs := os.Args
	os.Args = []string{"programm name", "test"}
	defer func() { os.Args = originArgs }()
	defer testhelper.CleanupTestEnvironment(t)
	testhelper.PreapreTestEnvironment(t)

	test := NewTestCommand()
	test.stdoutWriter = bytes.NewBufferString("")

	res := test.Test()

	expectedOutput := `
you must at least enter a commit message using parameter -m
Usage of git-commit-hook test:
  -b string
    	branch name
  -m string
    	commit message
  -p string
    	project name
`
	assert.Exactly(t, strings.TrimLeft(expectedOutput, "\n"), test.stdoutWriter.(*bytes.Buffer).String())
	assert.Exactly(t, 1, res)
}

func TestTestCommand_Test_InvalidParameters_PrintUsageAndErrorMessage(t *testing.T) {
	originArgs := os.Args
	os.Args = []string{"programm name", "test", "--unknown-parameter"}
	defer func() { os.Args = originArgs }()
	defer testhelper.CleanupTestEnvironment(t)
	testhelper.PreapreTestEnvironment(t)

	test := NewTestCommand()
	test.stdoutWriter = bytes.NewBufferString("")

	res := test.Test()

	expectedOutput := `
flag provided but not defined: -unknown-parameter
Usage of git-commit-hook test:
  -b string
    	branch name
  -m string
    	commit message
  -p string
    	project name
`
	assert.Exactly(t, strings.TrimLeft(expectedOutput, "\n"), test.stdoutWriter.(*bytes.Buffer).String())
	assert.Exactly(t, 1, res)
}

func TestTestCommand_Test_ConfigCannotBeFound(t *testing.T) {
	originArgs := os.Args
	os.Args = []string{"programm name", "test", ""}
	defer func() { os.Args = originArgs }()
	defer testhelper.CleanupTestEnvironment(t)
	testhelper.PreapreTestEnvironment(t)

	test := NewTestCommand()
	test.findConfigurationFilePath = func() (string, error) { return "", errors.New("configuration path not found") }
	test.stdoutWriter = bytes.NewBufferString("")

	res := test.Test()

	expectedOutput := `
error while searching configuration file: configuration path not found
`
	assert.Exactly(t, strings.TrimLeft(expectedOutput, "\n"), test.stdoutWriter.(*bytes.Buffer).String())
	assert.Exactly(t, 1, res)
}

func TestTestCommand_Test_LoadProjectConfiguration(t *testing.T) {
	originArgs := os.Args
	baseOsArgs := []string{"programm name", "test", "-m", "test commit", "-b", `"feature/PROJ-123"`}
	defer func() { os.Args = originArgs }()
	defer testhelper.CleanupTestEnvironment(t)
	testhelper.PreapreTestEnvironment(t)

	expectedProjectName := "sampleProject"

	testDataSet := map[string]struct {
		parameters             []string
		loadsConfigByParameter bool
	}{
		"no project parameter":   {parameters: []string{}, loadsConfigByParameter: false},
		"with project parameter": {parameters: []string{"-p", expectedProjectName}, loadsConfigByParameter: true},
	}

	for testCaseName, testData := range testDataSet {
		t.Run(testCaseName, func(t *testing.T) {
			os.Args = append(baseOsArgs, testData.parameters...)
			test := NewTestCommand()
			loadProjectConfigurationByNameCalled := false
			test.loadProjectConfigurationByName = func(p string) (config.Project, error) {
				if p != expectedProjectName {
					t.Fatalf("loadProjectConfigurationByName was expected to be called with p='%s', but got '%s'", expectedProjectName, p)
				}
				loadProjectConfigurationByNameCalled = true
				return config.Project{}, nil
			}
			loadProjectConfigurationFromWorkingDirCalled := false
			test.loadProjectConfigurationFromWorkingDir = func() (config.Project, error) {
				loadProjectConfigurationFromWorkingDirCalled = true
				return config.Project{}, nil
			}
			test.stdoutWriter = bytes.NewBufferString("")

			test.Test()

			if testData.loadsConfigByParameter {
				assert.True(t, loadProjectConfigurationByNameCalled)
				assert.False(t, loadProjectConfigurationFromWorkingDirCalled)
			} else {
				assert.False(t, loadProjectConfigurationByNameCalled)
				assert.True(t, loadProjectConfigurationFromWorkingDirCalled)
			}
		})
	}
}

func TestTestCommand_Test_LoadProjectConfigurationReturnsError_ExpectError(t *testing.T) {
	originArgs := os.Args
	os.Args = []string{"programm name", "test", "-m", "test commit"}
	defer func() { os.Args = originArgs }()
	defer testhelper.CleanupTestEnvironment(t)
	testhelper.PreapreTestEnvironment(t)

	expectedProjectName := "sampleProject"

	testDataSet := map[string]struct {
		parameters             []string
		loadsConfigByParameter bool
	}{
		"no project parameter":   {parameters: []string{}, loadsConfigByParameter: false},
		"with project parameter": {parameters: []string{"-p", expectedProjectName}, loadsConfigByParameter: true},
	}

	for testCaseName, testData := range testDataSet {
		t.Run(testCaseName, func(t *testing.T) {
			os.Args = append(os.Args, testData.parameters...)
			test := NewTestCommand()
			test.loadProjectConfigurationByName = func(string) (config.Project, error) {
				return config.Project{}, errors.New("some error")
			}
			test.loadProjectConfigurationFromWorkingDir = func() (config.Project, error) {
				return config.Project{}, errors.New("some error")
			}
			test.stdoutWriter = bytes.NewBufferString("")

			res := test.Test()
			assert.Exactly(t, 1, res)
		})
	}
}

func TestTestCommand_Test_BranchNameMissingAndWorkingDirIsNotGitRepo_ShowsError(t *testing.T) {
	originArgs := os.Args
	os.Args = []string{"programm name", "test", "-m", "test commit message", "-p", "test project"}
	defer func() { os.Args = originArgs }()
	defer testhelper.CleanupTestEnvironment(t)
	testhelper.PreapreTestEnvironment(t)

	test := NewTestCommand()
	test.stdoutWriter = bytes.NewBufferString("")

	res := test.Test()

	expectedOutput := `
error while reading branch name. ensure working dir is a git repo or use parameter -b to simulate a branch name
`
	assert.Exactly(t, strings.TrimLeft(expectedOutput, "\n"), test.stdoutWriter.(*bytes.Buffer).String())
	assert.Exactly(t, 1, res)
}

func TestTestCommand_Test_BranchNameMissingButWorkingDirIsGitRepo(t *testing.T) {
	originArgs := os.Args
	os.Args = []string{"programm name", "test", "-m", "test commit message", "-p", "test project"}
	defer func() { os.Args = originArgs }()
	defer testhelper.CleanupTestEnvironment(t)
	testhelper.PreapreTestEnvironment(t)
	testhelper.InitGitRepository(t, "feature/FROM-GIT")

	test := NewTestCommand()
	test.stdoutWriter = bytes.NewBufferString("")

	res := test.Test()

	expectedOutput := `
testing configuration '/tmp/git-commit-hook/git-commit-hook.yaml':
project        : test project
branch         : feature/FROM-GIT (current git branch)
commit message : test commit message

would generate the following commit message:
test commit message
`
	assert.Exactly(t, strings.TrimLeft(expectedOutput, "\n"), test.stdoutWriter.(*bytes.Buffer).String())
	assert.Exactly(t, 0, res)
}

func TestTestCommand_Test_ProjectNameNotFound_ShowsError(t *testing.T) {
	originArgs := os.Args
	defer func() { os.Args = originArgs }()
	defer testhelper.CleanupTestEnvironment(t)
	os.Args = []string{"programm name", "test", "-m", "test commit message", "-b", `"feature/PROJECT-123"`, "-p", "unknown project"}
	testhelper.PreapreTestEnvironment(t)
	testhelper.WriteConfigFile(t, "/tmp/git-commit-hook/")

	test := NewTestCommand()
	test.stdoutWriter = bytes.NewBufferString("")

	res := test.Test()

	expectedOutput := `
project configuration not found for project name 'unknown project'
`

	assert.Exactly(t, strings.TrimLeft(expectedOutput, "\n"), test.stdoutWriter.(*bytes.Buffer).String())
	assert.Exactly(t, 1, res)
}

func TestTestCommand_Test_CommitMessageModifierReturnsError_ShowError(t *testing.T) {
	originArgs := os.Args
	defer func() { os.Args = originArgs }()
	defer testhelper.CleanupTestEnvironment(t)
	os.Args = []string{"programm name", "test", "-m", "test commit message", "-b", `"feature/PROJECT-123"`, "-p", "test project"}
	testhelper.PreapreTestEnvironment(t)

	test := NewTestCommand()
	test.newCommitMessageModifier = func(project config.Project) hook.CommitMessageModifier {
		return &commitMessageModifierMock{}
	}
	test.stdoutWriter = bytes.NewBufferString("")

	res := test.Test()

	expectedOutput := `
testing configuration '/tmp/git-commit-hook/git-commit-hook.yaml':
project        : test project
branch         : "feature/PROJECT-123"
commit message : test commit message

some error while modifying the commit mesage
`
	assert.Exactly(t, strings.TrimLeft(expectedOutput, "\n"), test.stdoutWriter.(*bytes.Buffer).String())
	assert.Exactly(t, 1, res)
}

type commitMessageModifierMock struct{}

func (m *commitMessageModifierMock) ModifyGitCommitMessage(string, string) (string, error) {
	return "", errors.New("some error while modifying the commit mesage")
}
