package gitcommithook

import (
	"fmt"

	"github.com/Oppodelldog/git-commit-hook/git"
	"github.com/pkg/errors"
)

type gitBranchNameReaderFuncDef func() (string, error)
type featureBranchDetectFuncDef func(string) bool

var featureBranchDetectFunc = featureBranchDetectFuncDef(IsFeatureBranch)
var gitBranchNameReaderFunc = gitBranchNameReaderFuncDef(git.GetCurrentBranchName)

//ModifyGitCommitMessage prepends the current branch name to the given git commit message.
// if the current branch name is detected to be NO feature branch, the user will be prompted to enter
// a feature branch manually. This is then inserted in between current branch and commit message.
// If no valid branch name could be determined the function returns an error
func ModifyGitCommitMessage(gitCommitMessage string) (modifiedCommitMessage string, err error) {
	if gitCommitMessage == "" {
		err = errors.New("commit message is empty")
		return
	}

	branchName, err := gitBranchNameReaderFunc()
	if err != nil {
		return
	}

	if branchName == "" {
		err = errors.New("branch name was is empty")
		return
	}

	if !featureBranchDetectFunc(branchName) {
		if !featureBranchDetectFunc(gitCommitMessage) {
			err = errors.New("feature reference is required")
			return
		}
	}

	modifiedCommitMessage = fmt.Sprintf("%s: %s", branchName, gitCommitMessage)

	return modifiedCommitMessage, nil
}
