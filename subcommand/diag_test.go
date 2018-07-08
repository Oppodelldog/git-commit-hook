package subcommand

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/Oppodelldog/git-commit-hook/config"
	"github.com/Oppodelldog/git-commit-hook/testhelper"
	"github.com/stretchr/testify/assert"
)

func TestDiagCommand_Diagnostics(t *testing.T) {
	defer testhelper.CleanupTestEnvironment(t)
	testhelper.PreapreTestEnvironment(t)

	diag := NewDiagCommand()
	diag.stdoutWriter = bytes.NewBufferString("")

	res := diag.Diagnostics()

	expectedOutput := `
git-commit-hook diagnosticsload configuration: /tmp/git-commit-hook/git-commit-hook.yaml
-------------------------------------------------------------------
project:test projectpath   :/tmp/git-commit-hook/.git
branch types:
	feature:^feature/PROJECT-123$
	release:^release.*$

branch type templates:
	feature:{{.BranchName}}: {{.CommitMessage}}

branch type validation:
	release:
		(?m)(?:\s|^|/)(([A-Z](_)*)+-[0-9]+)([\s,;:!.-]|$):valid ticket ID

git-commit-hook installed: NO
`
	assert.Exactly(t, strings.TrimLeft(expectedOutput, "\n"), diag.stdoutWriter.(*bytes.Buffer).String())
	assert.Exactly(t, 0, res)
}

func TestDiagCommand_Diagnostics_CommitHookIsAlreadyInstalled(t *testing.T) {
	defer testhelper.CleanupTestEnvironment(t)

	testhelper.PreapreTestEnvironment(t)

	diag := NewDiagCommand()
	diag.stdoutWriter = bytes.NewBufferString("")
	diag.checkIsCommitHookInstalledAtPath = func(string) bool { return true }

	res := diag.Diagnostics()

	expectedOutput := `
git-commit-hook diagnosticsload configuration: /tmp/git-commit-hook/git-commit-hook.yaml
-------------------------------------------------------------------
project:test projectpath   :/tmp/git-commit-hook/.git
branch types:
	feature:^feature/PROJECT-123$
	release:^release.*$

branch type templates:
	feature:{{.BranchName}}: {{.CommitMessage}}

branch type validation:
	release:
		(?m)(?:\s|^|/)(([A-Z](_)*)+-[0-9]+)([\s,;:!.-]|$):valid ticket ID

git-commit-hook installed: YES
`
	assert.Exactly(t, strings.TrimLeft(expectedOutput, "\n"), diag.stdoutWriter.(*bytes.Buffer).String())
	assert.Exactly(t, 0, res)
}

func TestDiagCommand_Diagnostics_AnotherCommitHookIsAlreadyInstalled(t *testing.T) {
	defer testhelper.CleanupTestEnvironment(t)

	testhelper.PreapreTestEnvironment(t)

	diag := NewDiagCommand()
	diag.stdoutWriter = bytes.NewBufferString("")
	diag.checkIsCommitHookInstalledAtPath = func(string) bool { return false }
	diag.checkIsAnotherGitHookInstalledAtPath = func(string) bool { return true }

	res := diag.Diagnostics()

	expectedOutput := `
git-commit-hook diagnosticsload configuration: /tmp/git-commit-hook/git-commit-hook.yaml
-------------------------------------------------------------------
project:test projectpath   :/tmp/git-commit-hook/.git
branch types:
	feature:^feature/PROJECT-123$
	release:^release.*$

branch type templates:
	feature:{{.BranchName}}: {{.CommitMessage}}

branch type validation:
	release:
		(?m)(?:\s|^|/)(([A-Z](_)*)+-[0-9]+)([\s,;:!.-]|$):valid ticket ID

git-commit-hook installed: NO, another commit-msg hook is installed
`
	assert.Exactly(t, strings.TrimLeft(expectedOutput, "\n"), diag.stdoutWriter.(*bytes.Buffer).String())
	assert.Exactly(t, 0, res)
}

func TestDiagCommand_Diagnostics_ConfigFilePathCannotBeFound(t *testing.T) {
	defer testhelper.CleanupTestEnvironment(t)

	testhelper.PreapreTestEnvironment(t)

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
	defer testhelper.CleanupTestEnvironment(t)

	testhelper.PreapreTestEnvironment(t)

	diag := NewDiagCommand()
	diag.stdoutWriter = bytes.NewBufferString("")
	diag.loadConfiguration = func() (*config.Configuration, error) { return nil, errors.New("some error") }

	res := diag.Diagnostics()

	expectedOutput := `
git-commit-hook diagnosticsload configuration: /tmp/git-commit-hook/git-commit-hook.yaml
error loading configuration: some error
`
	assert.Exactly(t, strings.TrimLeft(expectedOutput, "\n"), diag.stdoutWriter.(*bytes.Buffer).String())
	assert.Exactly(t, 1, res)
}
