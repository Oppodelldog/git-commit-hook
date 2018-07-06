package hook

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCreateViewModel_BranchName(t *testing.T) {
	commitMessage := "\n\r\n\t\tHELLO\n\tWORLD\t\r\n\r\t"
	branchName := "456"
	viewModel := createViewModel(commitMessage, branchName)

	expectedViewModel := ViewModel{
		CommitMessage: "HELLO\n\tWORLD",
		BranchName:    branchName,
	}
	assert.Exactly(t, expectedViewModel, viewModel)
}
