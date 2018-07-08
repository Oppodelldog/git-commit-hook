package subcommand

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/Oppodelldog/git-commit-hook/config"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

const testDir = "/tmp/git-commit-hook/subcommand"

func TestDiagCommand_Diagnostics(t *testing.T) {
	defer cleanupTestEnvironment(t)

	preapreTestEnvironment(t)

	diag := NewDiagCommand()
	diag.stdoutWriter = bytes.NewBufferString("")

	res := diag.Diagnostics()

	expectedOutput := `
git-commit-hook diagnosticsload configuration: /tmp/git-commit-hook/subcommand/git-commit-hook.yaml
-------------------------------------------------------------------
project:test projectpath   :/tmp/git-commit-hook/subcommand/.git
branch types:	feature:^feature/PROJECT-123$	release:^release.*$
branch type templates:	feature:{{.BranchName}}: {{.CommitMessage}}
branch type validation:	release:		(?m)(?:\s|^|/)(([A-Z](_)*)+-[0-9]+)([\s,;:!.-]|$):valid ticket ID
git-commit-hook installed: NO
`
	assert.Exactly(t, strings.TrimLeft(expectedOutput, "\n"), diag.stdoutWriter.(*bytes.Buffer).String())
	assert.Exactly(t, 0, res)
}

func TestDiagCommand_Diagnostics_CommitHookIsAlreadyInstalled(t *testing.T) {
	defer cleanupTestEnvironment(t)

	preapreTestEnvironment(t)

	diag := NewDiagCommand()
	diag.stdoutWriter = bytes.NewBufferString("")
	diag.checkIsCommitHookInstalledAtPath = func(string) bool { return true }

	res := diag.Diagnostics()

	expectedOutput := `
git-commit-hook diagnosticsload configuration: /tmp/git-commit-hook/subcommand/git-commit-hook.yaml
-------------------------------------------------------------------
project:test projectpath   :/tmp/git-commit-hook/subcommand/.git
branch types:	feature:^feature/PROJECT-123$	release:^release.*$
branch type templates:	feature:{{.BranchName}}: {{.CommitMessage}}
branch type validation:	release:		(?m)(?:\s|^|/)(([A-Z](_)*)+-[0-9]+)([\s,;:!.-]|$):valid ticket ID
git-commit-hook installed: YES
`
	assert.Exactly(t, strings.TrimLeft(expectedOutput, "\n"), diag.stdoutWriter.(*bytes.Buffer).String())
	assert.Exactly(t, 0, res)
}

func TestDiagCommand_Diagnostics_AnotherCommitHookIsAlreadyInstalled(t *testing.T) {
	defer cleanupTestEnvironment(t)

	preapreTestEnvironment(t)

	diag := NewDiagCommand()
	diag.stdoutWriter = bytes.NewBufferString("")
	diag.checkIsCommitHookInstalledAtPath = func(string) bool { return false }
	diag.checkIsAnotherGitHookInstalledAtPath = func(string) bool { return true }

	res := diag.Diagnostics()

	expectedOutput := `
git-commit-hook diagnosticsload configuration: /tmp/git-commit-hook/subcommand/git-commit-hook.yaml
-------------------------------------------------------------------
project:test projectpath   :/tmp/git-commit-hook/subcommand/.git
branch types:	feature:^feature/PROJECT-123$	release:^release.*$
branch type templates:	feature:{{.BranchName}}: {{.CommitMessage}}
branch type validation:	release:		(?m)(?:\s|^|/)(([A-Z](_)*)+-[0-9]+)([\s,;:!.-]|$):valid ticket ID
git-commit-hook installed: NO, another commit-msg hook is installed
`
	assert.Exactly(t, strings.TrimLeft(expectedOutput, "\n"), diag.stdoutWriter.(*bytes.Buffer).String())
	assert.Exactly(t, 0, res)
}

func TestDiagCommand_Diagnostics_ConfigFilePathCannotBeFound(t *testing.T) {
	defer cleanupTestEnvironment(t)

	preapreTestEnvironment(t)

	diag := NewDiagCommand()
	diag.stdoutWriter = bytes.NewBufferString("")
	diag.findConfigurationFilePath = func() (string, error) { return "", errors.New("some error") }

	res := diag.Diagnostics()

	expectedOutput := `
error while searching configuration file: some error
`
	assert.Exactly(t, strings.TrimLeft(expectedOutput, "\n"), diag.stdoutWriter.(*bytes.Buffer).String())
	assert.Exactly(t, 1, res)
}

func TestDiagCommand_Diagnostics_ConfigCannotBeLoad(t *testing.T) {
	defer cleanupTestEnvironment(t)

	preapreTestEnvironment(t)

	diag := NewDiagCommand()
	diag.stdoutWriter = bytes.NewBufferString("")
	diag.loadConfiguration = func() (*config.Configuration, error) { return nil, errors.New("some error") }

	res := diag.Diagnostics()

	expectedOutput := `
git-commit-hook diagnosticsload configuration: /tmp/git-commit-hook/subcommand/git-commit-hook.yaml
error loading configuration: some error
`
	assert.Exactly(t, strings.TrimLeft(expectedOutput, "\n"), diag.stdoutWriter.(*bytes.Buffer).String())
	assert.Exactly(t, 1, res)
}

func cleanupTestEnvironment(t *testing.T) {
	err := os.RemoveAll(testDir)
	if err != nil {
		t.Fatalf("Error cleaning up test environment.. Did not expect os.RemoveAll to return an error, but got: %v ", err)
	}
}

func preapreTestEnvironment(t *testing.T) {
	writeConfigFile(t, testDir)
	err := os.Chdir(testDir)
	if err != nil {
		t.Fatalf("Error preparing test environment.Did not expect os.Chdir to return an error, but got: %v ", err)
	}
}

func writeConfigFile(t *testing.T, dir string) {
	os.MkdirAll(dir, 0777)
	cfg := config.Configuration{
		"test project": config.Project{
			Path: "/tmp/git-commit-hook/subcommand/.git",
			BranchTypes: map[string]config.BranchTypePattern{
				"feature": `^feature/PROJECT-123$`,
				"release": `^release.*$`,
			},
			Templates: map[string]config.BranchTypeTemplate{
				"feature": "{{.BranchName}}: {{.CommitMessage}}",
			},
			Validation: map[string]config.BranchValidationConfiguration{
				"release": {
					"(?m)(?:\\s|^|/)(([A-Z](_)*)+-[0-9]+)([\\s,;:!.-]|$)": "valid ticket ID",
				},
			},
		},
	}

	configBytes, err := yaml.Marshal(cfg)
	if err != nil {
		t.Fatalf("Did not expect yaml.Marshal to return an error, but got: %v ", err)
	}

	err = ioutil.WriteFile(path.Join(dir, "git-commit-hook.yaml"), configBytes, 0666)
	if err != nil {
		t.Fatalf("Did not expect ioutil.WriteFile to return an error, but got: %v ", err)
	}
}
