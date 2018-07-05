package main

import (
	"os"

	"path/filepath"

	"fmt"

	"github.com/Oppodelldog/git-commit-hook"
	"github.com/Oppodelldog/git-commit-hook/config"
	"github.com/pkg/errors"
)

type rewriteCommitMessageFuncDef func(string, config.ProjectConfiguration) error
type exitFuncDef func(code int)

var rewriteCommitMessageFunc = rewriteCommitMessageFuncDef(gitcommithook.RewriteCommitMessage)
var exitFunc = exitFuncDef(os.Exit)

func main() {
	if len(os.Args) < 2 {
		fmt.Print(errors.New("no input provided"))
		exitFunc(1)
		return
	}

	commitMessageFile := os.Args[1]
	if commitMessageFile == "" {
		fmt.Print(errors.New("no commit message file passed as parameter 1"))
		exitFunc(1)
		return
	}

	projectConfiguration, err := loadProjectConfiguration(commitMessageFile)
	if err != nil {
		fmt.Print(err)
		exitFunc(1)
		return
	}

	err = rewriteCommitMessageFunc(commitMessageFile, projectConfiguration)
	if err != nil {
		fmt.Print(err)
		exitFunc(1)
		return
	}

	exitFunc(0)
}

func loadProjectConfiguration(commitMessageFile string) (config.ProjectConfiguration, error) {

	configuration, err := config.LoadConfiguration()
	if err != nil {
		return config.ProjectConfiguration{}, err
	}

	projectPath, err := filepath.Abs(filepath.Dir(commitMessageFile))
	if err != nil {
		return config.ProjectConfiguration{}, err
	}

	return configuration.GetProjectConfiguration(projectPath)
}
