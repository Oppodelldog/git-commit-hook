package hook

import (
	"github.com/Oppodelldog/git-commit-hook/config"
	"github.com/pkg/errors"
)

type (
	// CommitMessageModifier implements the modification of a given commit message and the current branch name.
	CommitMessageModifier interface {
		ModifyGitCommitMessage(gitCommitMessage, branchName string) (modifiedCommitMessage string, err error)
	}
	commitMessageModifier struct {
		createViewModelFunc       createViewModelFuncDef
		renderCommitMessageFunc   renderCommitMessageFuncDef
		validateCommitmessageFunc validateCommitMessageFuncDef
	}

	createViewModelFuncDef       func(gitCommitMessage string, branchName string) ViewModel
	validateCommitMessageFuncDef func(branchName, modifiedCommitMessage string) error
	renderCommitMessageFuncDef   func(branchName string, viewModel ViewModel) (string, error)
)

// NewCommitMessageModifier create a CommitMessageModifier
func NewCommitMessageModifier(projectConfiguration config.Project) CommitMessageModifier {
	commitMessageRenderer := &CommitMessageRenderer{projectConfiguration}
	return &commitMessageModifier{
		createViewModelFunc:       createViewModel,
		renderCommitMessageFunc:   commitMessageRenderer.RenderCommitMessage,
		validateCommitmessageFunc: projectConfiguration.Validate,
	}
}

//ModifyGitCommitMessage prepends the current branch name to the given git commit message.
// if the current branch name is detected to be NO feature branch, the user will be prompted to enter
// a feature branch manually. This is then inserted in between current branch and commit message.
// If no valid branch name could be determined the function returns an error
func (m *commitMessageModifier) ModifyGitCommitMessage(gitCommitMessage string, branchName string) (modifiedCommitMessage string, err error) {

	modifiedCommitMessage = gitCommitMessage
	if gitCommitMessage == "" {
		err = errors.New("commit message is empty")
		return
	}

	if branchName == "" {
		return
	}

	viewModel := createViewModel(gitCommitMessage, branchName)

	modifiedCommitMessage, err = m.renderCommitMessageFunc(branchName, viewModel)
	if err != nil {
		return
	}

	err = m.validateCommitmessageFunc(branchName, modifiedCommitMessage)
	if err != nil {
		modifiedCommitMessage = ""
	}

	return
}
