package gitcommithook

import (
	"os"
	"testing"

	"github.com/Oppodelldog/git-commit-hook/config"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

const commitMessageFile = ".git/COMMIT_MSG_EDIT"
const commitMessage = "commitMessage"
const modifiedCommitMessage = "modified commitMessage"

var rewriteOriginals = struct {
	readFileFunc  readFileFuncDef
	modifyFunc    modifyFuncDef
	writeFileFunc writeFileFuncDef
}{
	readFileFunc:  readFileFunc,
	modifyFunc:    modifyFunc,
	writeFileFunc: writeFileFunc,
}

func restoreRewriteOriginals() {
	readFileFunc = rewriteOriginals.readFileFunc
	modifyFunc = rewriteOriginals.modifyFunc
	writeFileFunc = rewriteOriginals.writeFileFunc
}

func TestRewriteCommitMessage_ErrorCase_CannotFindCommitMessageFile(t *testing.T) {
	defer restoreRewriteOriginals()

	commitMessageFileCannotBeFound(t)

	err := RewriteCommitMessage(commitMessageFile, config.ProjectConfiguration{})

	assert.Contains(t, err.Error(), "error reading commit message")
}

func TestRewriteCommitMessage_ErrorCase_ErrorModifyingMessage(t *testing.T) {
	defer restoreRewriteOriginals()
	commitMessageFileCanBeRead(t)
	errorModifyingCommitMessage(t)

	err := RewriteCommitMessage(commitMessageFile, config.ProjectConfiguration{})

	assert.Contains(t, err.Error(), "error modifying commit message")
}

func TestRewriteCommitMessage_ErrorCase_CannotWriteMessageToFile(t *testing.T) {
	defer restoreRewriteOriginals()

	commitMessageFileCanBeRead(t)
	commitMessageIsModified(t)
	cannotWriteToFile(t)

	err := RewriteCommitMessage(commitMessageFile, config.ProjectConfiguration{})

	assert.Contains(t, err.Error(), "error writing commit message to")
}

func TestRewriteCommitMessage_HappyPath(t *testing.T) {
	defer restoreRewriteOriginals()

	commitMessageFileCanBeRead(t)
	commitMessageIsModified(t)
	commitMessageIsWrittenToFile(t)

	err := RewriteCommitMessage(commitMessageFile, config.ProjectConfiguration{})

	assert.NoError(t, err)
}

func cannotWriteToFile(t *testing.T) {
	writeFileFunc = func(fileName string, data []byte, perm os.FileMode) error {
		assert.Exactly(t, commitMessageFile, fileName)
		assert.Exactly(t, modifiedCommitMessage, string(data))
		assert.Exactly(t, perm, os.FileMode(0777))

		return errors.New("cannot write to file")
	}
}

func commitMessageIsWrittenToFile(t *testing.T) {
	writeFileFunc = func(fileName string, data []byte, perm os.FileMode) error {
		assert.Exactly(t, commitMessageFile, fileName)
		assert.Exactly(t, modifiedCommitMessage, string(data))
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

func errorModifyingCommitMessage(t *testing.T) {
	modifyFunc = func(message string, configuration config.ProjectConfiguration) (string, error) {
		assert.Exactly(t, commitMessage, message)
		return "", errors.New("error modifying")
	}
}

func commitMessageIsModified(t *testing.T) {
	modifyFunc = func(message string, configuration config.ProjectConfiguration) (string, error) {
		assert.Exactly(t, commitMessage, message)
		return modifiedCommitMessage, nil
	}
}
