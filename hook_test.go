package gitcommithook

import (
	"testing"
	"github.com/magiconair/properties/assert"
)

func TestModifyGitCommitMessage(t *testing.T) {
	inputGitCommitMessage := "initial commit"

	modifiedGitCommitMessage := ModifyGitCommitMessage(inputGitCommitMessage)

	expectedGitCommitMessage := "BRANCH: initial commit"

	assert.Equal(t, expectedGitCommitMessage, modifiedGitCommitMessage)
}
