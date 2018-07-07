package hook

import (
	"os"
	"testing"

	"github.com/Oppodelldog/git-commit-hook/config"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

const commitMessageFile = ".git/COMMIT_MSG_EDIT"
const commitMessage = "commitMessage"

var rewriteOriginals = struct {
	readFileFunc  readFileFuncDef
	writeFileFunc writeFileFuncDef
}{
	readFileFunc:  readFileFunc,
	writeFileFunc: writeFileFunc,
}

func restoreRewriteOriginals() {
	readFileFunc = rewriteOriginals.readFileFunc
	writeFileFunc = rewriteOriginals.writeFileFunc
}

func TestRewriteCommitMessage_ErrorCase_CannotFindCommitMessageFile(t *testing.T) {
	defer restoreRewriteOriginals()

	commitMessageFileCannotBeFound(t)
	modifier := NewCommitMessageModifier(config.Project{})

	err := RewriteCommitMessage(commitMessageFile, modifier)

	assert.Contains(t, err.Error(), "error reading commit message")
}

func TestRewriteCommitMessage_ErrorCase_ErrorModifyingMessage(t *testing.T) {
	defer restoreRewriteOriginals()
	commitMessageFileCanBeRead(t)
	modifier := &commitMessageModifierStub{"", errors.New("some error")}

	err := RewriteCommitMessage(commitMessageFile, modifier)

	assert.Contains(t, err.Error(), "some error")
}

type commitMessageModifierStub struct {
	gitCommitMessage string
	err              error
}

func (m *commitMessageModifierStub) ModifyGitCommitMessage(gitCommitMessage string, branchName string) (modifiedCommitMessage string, err error) {
	return m.gitCommitMessage, m.err
}

func TestRewriteCommitMessage_ErrorCase_CannotWriteMessageToFile(t *testing.T) {
	defer restoreRewriteOriginals()

	commitMessageFileCanBeRead(t)
	cannotWriteToFile(t)
	modifier := NewCommitMessageModifier(config.Project{})

	err := RewriteCommitMessage(commitMessageFile, modifier)

	assert.Contains(t, err.Error(), "error writing commit message to")
}

func TestRewriteCommitMessage_HappyPath(t *testing.T) {
	defer restoreRewriteOriginals()

	commitMessageFileCanBeRead(t)
	commitMessageIsWrittenToFile(t)
	modifier := NewCommitMessageModifier(config.Project{})

	err := RewriteCommitMessage(commitMessageFile, modifier)

	assert.NoError(t, err)
}

func cannotWriteToFile(t *testing.T) {
	writeFileFunc = func(fileName string, data []byte, perm os.FileMode) error {
		assert.Exactly(t, commitMessageFile, fileName)
		assert.Exactly(t, commitMessage, string(data))
		assert.Exactly(t, perm, os.FileMode(0777))

		return errors.New("cannot write to file")
	}
}

func commitMessageIsWrittenToFile(t *testing.T) {
	writeFileFunc = func(fileName string, data []byte, perm os.FileMode) error {
		assert.Exactly(t, commitMessageFile, fileName)
		assert.Exactly(t, commitMessage, string(data))
		assert.Exactly(t, perm, os.FileMode(0777))

		return nil
	}
}

func commitMessageFileCanBeRead(t *testing.T) {
	readFileFunc = func(fileName string) ([]byte, error) {
		assert.Exactly(t, commitMessageFile, fileName)
		return []byte(commitMessage), nil
	}
}

func commitMessageFileCannotBeFound(t *testing.T) {
	readFileFunc = func(fileName string) ([]byte, error) {
		assert.Exactly(t, commitMessageFile, fileName)
		return []byte{}, errors.New("Could not find file")
	}
}
