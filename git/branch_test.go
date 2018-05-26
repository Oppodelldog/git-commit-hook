package git

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

var originals = struct {
	execFunc execFuncDef
}{
	execFunc: execFunc,
}

func restoreOriginals() {
	execFunc = originals.execFunc
}

func TestGetCurrentBranchName_FindsBranchInCommandOutput(t *testing.T) {
	defer restoreOriginals()

	testCases := []struct {
		gitOutput          string
		expectedBranchName string
	}{
		{``, ""},
		{`develop`, ""},
		{`* develop`, "develop"},
		{"* develop\nmaster", "develop"},
		{"release/v1.0.1\n* develop\nmaster", "develop"},
		{"feature/xyz\ndevelop\n* release/v1.0.1\nmaster", "release/v1.0.1"},
	}

	for i, testCase := range testCases {
		t.Run(string(i), func(t *testing.T) {
			execFunc = func(s1 string, s2 ...string) *exec.Cmd {
				return exec.Command("echo", testCase.gitOutput)
			}
			branchName, err := GetCurrentBranchName()
			assert.NoError(t, err)
			assert.Exactly(t, testCase.expectedBranchName, branchName)
		})
	}
}

func TestGetCurrentBranchName_ExpectGitCommandIsCalledProperly(t *testing.T) {
	defer restoreOriginals()

	execFunc = func(s1 string, s2 ...string) *exec.Cmd {
		assert.Exactly(t, "git", s1)
		assert.Exactly(t, []string{"branch"}, s2)

		return exec.Command("", "")
	}

	GetCurrentBranchName()
}

func TestGetCurrentBranchName_ReturnsErrorFromGitCommandExecution(t *testing.T) {
	defer restoreOriginals()

	execFunc = func(s1 string, s2 ...string) *exec.Cmd {
		//noinspection SpellCheckingInspection
		return exec.Command("thiscommandwillnotbefound", "")
	}

	_, err := GetCurrentBranchName()
	assert.Error(t, err)
}

func TestDefaultExecFuncIsExecCommand(t *testing.T) {
	assert.IsType(t, execFuncDef(exec.Command), execFunc)
}
