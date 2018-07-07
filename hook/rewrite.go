package hook

import (
	"io/ioutil"

	"os"

	"github.com/Oppodelldog/git-commit-hook/config"
	"github.com/pkg/errors"
	"github.com/Oppodelldog/git-commit-hook/git"
)

type (
	readFileFuncDef func(filename string) ([]byte, error)
	modifyFuncDef func(string, config.Project) (string, error)
	writeFileFuncDef func(string, []byte, os.FileMode) error
)

var (
	readFileFunc  = readFileFuncDef(ioutil.ReadFile)
	writeFileFunc = writeFileFuncDef(ioutil.WriteFile)
)

//RewriteCommitMessage rewrites the commit message in the given commit message file
func RewriteCommitMessage(commitMessageFile string, commitMessageModifier CommitMessageModifier) error {

	var outputMessage string
	fileContent, err := readFileFunc(commitMessageFile)
	if err != nil {
		return errors.Errorf("error reading commit message from '%s': %v", commitMessageFile, err.Error())
	}

	branchName, _ := git.GetCurrentBranchName()
	outputMessage, err = commitMessageModifier.ModifyGitCommitMessage(string(fileContent), branchName)
	if err != nil {
		return errors.Errorf("error modifying commit message: %s", err.Error())
	}

	err = writeFileFunc(commitMessageFile, []byte(outputMessage), 0777)
	if err != nil {
		return errors.Errorf("error writing commit message to '%s': %s", commitMessageFile, err.Error())
	}

	return nil
}
