package gitcommithook

import (
	"testing"

	"github.com/Oppodelldog/git-commit-hook/git"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

//originals hold a compile time copy of the modules globals
//this enables stubbing functions or user input
var originals = struct {
	gitBranchNameReaderFunc gitBranchNameReaderFuncDef
	featureBranchDetectFunc featureBranchDetectFuncDef
}{
	gitBranchNameReaderFunc: gitBranchNameReaderFunc,
	featureBranchDetectFunc: featureBranchDetectFunc,
}

//restoreOriginals resets the modules globals to the original state
func restoreOriginals() {
	gitBranchNameReaderFunc = originals.gitBranchNameReaderFunc
	featureBranchDetectFunc = originals.featureBranchDetectFunc
}

func TestModifyGitCommitMessage(t *testing.T) {

	testCases := map[string]struct {
		input        string
		branchName   string
		output       string
		returnsError bool
	}{
		"both set": {
			input:        "initial commit",
			branchName:   "BRANCH",
			output:       "BRANCH: initial commit",
			returnsError: false,
		},
		"empty input": {
			input:        "",
			branchName:   "BRANCH",
			output:       "",
			returnsError: true,
		},
		"empty branchName": {
			input:        "initial commit",
			branchName:   "",
			output:       "",
			returnsError: true,
		},
		"both empty": {
			input:        "",
			branchName:   "",
			output:       "",
			returnsError: true,
		},
	}

	for testName, testData := range testCases {

		t.Run(testName, func(t *testing.T) {
			branchNameWillBe(testData.branchName)

			modifiedGitCommitMessage, err := ModifyGitCommitMessage(testData.input)

			if testData.returnsError {
				assert.Error(t, err)
			}
			assert.Equal(t, testData.output, modifiedGitCommitMessage)

			restoreOriginals()
		})
	}
}

func TestModifyGitCommitMessage_branchNameReaderReturnsError_ExpectError(t *testing.T) {
	defer restoreOriginals()

	expectedError := errors.New("error while reading branch name")

	branchNameReaderReturnsError(expectedError)

	_, err := ModifyGitCommitMessage("some message")

	assert.Equal(t, expectedError, err)
}

func TestModifyGitCommitMessage_FixesRequireFeatureBranchReference(t *testing.T) {
	defer restoreOriginals()

	branchNameWillBe("no-feature-branch-name")
	branchWillBeNoFeatureBranch()

	inputMessage := "some rc-fix, what was the feature again? Hmm, I leave it out."

	modifiedCommitMessage, err := ModifyGitCommitMessage(inputMessage)

	assert.Error(t, err)
	assert.Equal(t, "", modifiedCommitMessage)
}

func TestDefaultGitBranchNameReaderFunc(t *testing.T) {
	assert.IsType(t, gitBranchNameReaderFuncDef(git.GetCurrentBranchName), gitBranchNameReaderFunc)
}

func TestDefaultFeatureBranchDetectFunc(t *testing.T) {
	assert.IsType(t, featureBranchDetectFuncDef(IsFeatureBranch), featureBranchDetectFunc)
}

func branchWillBeNoFeatureBranch() {
	featureBranchDetectFunc = func(branchName string) bool { return false }
}

func branchNameWillBe(s string) {
	gitBranchNameReaderFunc = func() (string, error) { return s, nil }
}

func branchNameReaderReturnsError(e error) {
	gitBranchNameReaderFunc = func() (string, error) { return "", e }
}
