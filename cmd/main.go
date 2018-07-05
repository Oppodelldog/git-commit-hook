package main

import (
	"os"

	"fmt"

	"github.com/Oppodelldog/git-commit-hook"
	"github.com/Oppodelldog/git-commit-hook/config"
	"github.com/pkg/errors"
	"github.com/Oppodelldog/git-commit-hook/cmd/diag"
)

type diagnosticsFuncDef func() int
type rewriteCommitMessageFuncDef func(string, config.ProjectConfiguration) error
type exitFuncDef func(code int)

var diagnosticsFunc = diag.SubCommandDiagnostics
var rewriteCommitMessageFunc = rewriteCommitMessageFuncDef(gitcommithook.RewriteCommitMessage)
var exitFunc = exitFuncDef(os.Exit)

func main() {
	if len(os.Args) < 2 {
		result := diagnosticsFunc()
		exitFunc(result)
		return
	}

	commitMessageFile := os.Args[1]
	if commitMessageFile == "" {
		fmt.Print(errors.New("no commit message file passed as parameter 1"))
		exitFunc(1)
		return
	}

	projectConfiguration, err := config.LoadProjectConfiguration(commitMessageFile)
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
