package hook

import (
	"strings"

	"github.com/Oppodelldog/git-commit-hook/config"
	"github.com/Oppodelldog/git-commit-hook/git"
	"github.com/pkg/errors"
)

type gitBranchNameReaderFuncDef func() (string, error)

var gitBranchNameReaderFunc = gitBranchNameReaderFuncDef(git.GetCurrentBranchName)

//ModifyGitCommitMessage prepends the current branch name to the given git commit message.
// if the current branch name is detected to be NO feature branch, the user will be prompted to enter
// a feature branch manually. This is then inserted in between current branch and commit message.
// If no valid branch name could be determined the function returns an error
func ModifyGitCommitMessage(gitCommitMessage string, projectConfiguration config.ProjectConfiguration) (modifiedCommitMessage string, err error) {

	if gitCommitMessage == "" {
		err = errors.New("commit message is empty")
		return
	}

	branchName, err := gitBranchNameReaderFunc()
	if err != nil {
		return
	}

	if branchName == "" {
		err = errors.New("branch name is empty")
		return
	}

	viewModel := createViewModel(gitCommitMessage, branchName)
	modifiedCommitMessage, err = projectConfiguration.RenderCommitMessage(branchName, viewModel)
	if err != nil {
		return
	}

	err = projectConfiguration.Validate(branchName, modifiedCommitMessage)
	if err != nil {
		modifiedCommitMessage = ""
	}

	return
}

func ModifyGitCommitMessageForCustomBranch(gitCommitMessage string, projectConfiguration config.ProjectConfiguration, branchNameResolver gitBranchNameReaderFuncDef) (modifiedCommitMessage string, err error) {
	originalGitBranchNameReaderFunc := gitBranchNameReaderFunc
	defer func() { gitBranchNameReaderFunc = originalGitBranchNameReaderFunc }()
	gitBranchNameReaderFunc = branchNameResolver

	return ModifyGitCommitMessage(gitCommitMessage, projectConfiguration)
}

func createViewModel(gitCommitMessage string, branchName string) config.ViewModel {
	trimmedCommitMessage := strings.Trim(gitCommitMessage, " \t\r\n")
	viewModel := config.ViewModel{CommitMessage: trimmedCommitMessage, BranchName: branchName}

	return viewModel
}
