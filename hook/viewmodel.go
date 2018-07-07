package hook

import (
	"strings"
)

// ViewModel defines all variables that can be in templates to define the modified commit message
type (
	ViewModel struct {
		BranchName    string
		CommitMessage string
	}
)

func createViewModel(commitMessage string, branchName string) ViewModel {
	trimmedCommitMessage := strings.Trim(commitMessage, " \t\r\n")
	viewModel := ViewModel{CommitMessage: trimmedCommitMessage, BranchName: branchName}

	return viewModel
}
