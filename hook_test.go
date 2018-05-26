package gitcommithook

import (
	"testing"

	"reflect"

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
		input        string
		branchName   string
		output       string
		returnsError bool
	}{
		"commit for feature branch": {
			input:        "initial commit",
			branchName:   "PROJECT-123",
			output:       "PROJECT-123: initial commit",
			returnsError: false,
		},
		"commit for fix-branch with no feature reference": {
			input:        "initial commit",
			branchName:   "release/v1.0.1-fix",
			output:       "",
			returnsError: true,
		},
		"commit for fix-branch with feature reference": {
			input:        "feature/PROJECT-123 initial commit",
			branchName:   "release/v1.0.1-fix",
			output:       "release/v1.0.1-fix: feature/PROJECT-123 initial commit",
			returnsError: false,
		},
		"commit for fix-branch with feature reference somewhere in the commit message": {
			input:        "fixed something for PROJECT-123, should work now",
			branchName:   "release/v1.0.1-fix",
			output:       "release/v1.0.1-fix: fixed something for PROJECT-123, should work now",
			returnsError: false,
		},
		"commit for feature without commit message": {
			input:        "",
			branchName:   "PROJECT-123",
			output:       "",
			returnsError: true,
		},
		"commit with broken branch name detection, but with commit message": {
			input:        "initial commit",
			branchName:   "",
			output:       "",
			returnsError: true,
		},
		"commit with broken branch name detection and without commit message": {
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

	_, err := ModifyGitCommitMessage("some message")

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
