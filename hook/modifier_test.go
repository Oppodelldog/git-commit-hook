package hook

import (
	"testing"

	"github.com/Oppodelldog/git-commit-hook/config"
	"github.com/stretchr/testify/assert"
	"errors"
)

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
			output:        "initial commit",
			errorContains: "",
		},
		"commit with broken branch name detection and without commit message": {
			input:         "",
			branchName:    "",
			output:        "",
			errorContains: "commit message is empty",
		},
	}

	prjCfg := config.Project{
		BranchTypes: map[string]config.BranchTypePattern{
			"feature": `(?m)^((?!master|release|develop).)*$`,
			"release": `(?m)^(origin\/)*release\/v([0-9]*\.*)*(-fix)*$`,
		},
		Templates: map[string]config.BranchTypeTemplate{
			"feature": "{{.BranchName}}: {{.CommitMessage}}",
			"release": "{{.BranchName}}: {{.CommitMessage}}",
		},
		Validation: map[string]config.BranchValidationConfiguration{
			"release": {
				`(?m)(?:\s|^|/)(([A-Z](_)*)+-[0-9]+)([\s,;:!.-]|$)`: "valid feature reference ID",
			},
		},
	}

	for testName, testData := range testCases {

		t.Run(testName, func(t *testing.T) {
			modifier := NewCommitMessageModifier(prjCfg)
			modifiedGitCommitMessage, err := modifier.ModifyGitCommitMessage(testData.input, testData.branchName)

			if testData.errorContains != "" {
				assert.Contains(t, err.Error(), testData.errorContains)
			} else {
				assert.NoError(t, err)
			}
			assert.Exactly(t, testData.output, modifiedGitCommitMessage)
		})
	}
}

func TestModifyGitCommitMessage_commitMessageRendererReturnsError_ExpectError(t *testing.T) {

	commitMessage := "some message"
	branchName := "feature123"
	errStub := errors.New("stubbed renderer error")
	modifier := NewCommitMessageModifier(config.Project{})
	modifier.(*commitMessageModifier).renderCommitMessageFunc = func(branchName string, viewModel ViewModel) (string, error) {
		return "", errStub
	}
	_, err := modifier.ModifyGitCommitMessage(commitMessage, branchName)

	assert.Exactly(t, errStub, err)
}
