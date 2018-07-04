package gitcommithook

import (
	"testing"

	"reflect"

	"github.com/Oppodelldog/git-commit-hook/config"
	"github.com/Oppodelldog/git-commit-hook/git"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

//originals hold a compile time copy of the modules globals
//this enables stubbing functions or user input
var originals = struct {
	gitBranchNameReaderFunc gitBranchNameReaderFuncDef
}{
	gitBranchNameReaderFunc: gitBranchNameReaderFunc,
}

//restoreOriginals resets the modules globals to the original state
func restoreOriginals() {
	gitBranchNameReaderFunc = originals.gitBranchNameReaderFunc
}

func TestModifyGitCommitMessage(t *testing.T) {

	testCases := map[string]struct {
		input         string
		branchName    string
		output        string
		errorContains string
	}{
		"commit for feature branch": {
			input:      "initial commit",
			branchName: "PROJECT-123",
			output:     "PROJECT-123: initial commit",
		},
		"commit for fix-branch with no feature reference": {
			input:         "initial commit",
			branchName:    "release/v1.0.1-fix",
			output:        "",
			errorContains: "validation error for branch 'release/v1.0.1-fix'",
		},
		"commit for fix-branch with feature reference": {
			input:      "feature/PROJECT-123 initial commit",
			branchName: "release/v1.0.1-fix",
			output:     "release/v1.0.1-fix: feature/PROJECT-123 initial commit",
		},
		"commit for fix-branch with feature reference somewhere in the commit message": {
			input:      "fixed something for PROJECT-123, should work now",
			branchName: "release/v1.0.1-fix",
			output:     "release/v1.0.1-fix: fixed something for PROJECT-123, should work now",
		},
		"commit for feature without commit message": {
			input:         "",
			branchName:    "PROJECT-123",
			output:        "",
			errorContains: "commit message is empty",
		},
		"commit with broken branch name detection, but with commit message": {
			input:         "initial commit",
			branchName:    "",
			output:        "",
			errorContains: "branch name is empty",
		},
		"commit with broken branch name detection and without commit message": {
			input:         "",
			branchName:    "",
			output:        "",
			errorContains: "commit message is empty",
		},
	}

	prjCfg := config.ProjectConfiguration{
		BranchTypes: map[string]config.BranchTypeConfiguration{
			"feature": {Pattern: `(?m)^((?!master|release|develop).)*$`},
			"release": {Pattern: `(?m)^(origin\/)*release\/v([0-9]*\.*)*(-fix)*$`},
		},
		Templates: map[string]config.BranchTemplateConfiguration{
			"feature": {Template: "{{.BranchName}}: {{.CommitMessage}}"},
			"release": {Template: "{{.BranchName}}: {{.CommitMessage}}"},
		},
		Validation: map[string]config.BranchValidationConfiguration{
			"release": {
				`(?m)(?:\s|^|/)(([A-Z](_)*)+-[0-9]+)([\s,;:!.-]|$)`: "valid feature reference ID",
			},
		},
	}

	for testName, testData := range testCases {

		t.Run(testName, func(t *testing.T) {
			branchNameWillBe(testData.branchName)

			modifiedGitCommitMessage, err := ModifyGitCommitMessage(testData.input, prjCfg)

			if testData.errorContains != "" {
				assert.Contains(t, err.Error(), testData.errorContains)
			} else {
				assert.NoError(t, err)
			}
			assert.Exactly(t, testData.output, modifiedGitCommitMessage)

			restoreOriginals()
		})
	}
}

func TestModifyGitCommitMessage_branchNameReaderReturnsError_ExpectError(t *testing.T) {
	defer restoreOriginals()

	expectedError := errors.New("error while reading branch name")

	branchNameReaderReturnsError(expectedError)

	_, err := ModifyGitCommitMessage("some message", config.ProjectConfiguration{})

	assert.Exactly(t, expectedError, err)
}

func TestDefaultGitBranchNameReaderFunc(t *testing.T) {
	assert.Exactly(t, reflect.ValueOf(git.GetCurrentBranchName).Pointer(), reflect.ValueOf(gitBranchNameReaderFunc).Pointer())
}

func TestDefaultFeatureBranchDetectFunc(t *testing.T) {
	assert.Exactly(t, reflect.ValueOf(IsFeatureBranch).Pointer(), reflect.ValueOf(featureBranchDetectFunc).Pointer())
}

func branchNameWillBe(s string) {
	gitBranchNameReaderFunc = func() (string, error) { return s, nil }
}

func branchNameReaderReturnsError(e error) {
	gitBranchNameReaderFunc = func() (string, error) { return "", e }
}
