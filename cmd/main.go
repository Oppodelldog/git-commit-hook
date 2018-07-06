package main

import (
	"os"

	"fmt"

	"github.com/Oppodelldog/git-commit-hook/cmd/diag"
	"github.com/Oppodelldog/git-commit-hook/config"
	"github.com/Oppodelldog/git-commit-hook/hook"
	"github.com/pkg/errors"
)

type callWithIntResult func() int
type rewriteCommitMessageFuncDef func(string, config.ProjectConfiguration) error
type exitFuncDef func(code int)

var diagnosticsFunc = diag.SubCommandDiagnostics
var testFunc = diag.Test
var installFunc = diag.Install
var uninstallFunc = diag.Uninstall
var rewriteCommitMessageFunc = rewriteCommitMessageFuncDef(hook.RewriteCommitMessage)
var exitFunc = exitFuncDef(os.Exit)

func main() {
	if len(os.Args) < 2 {
		result := diagnosticsFunc()
		exitFunc(result)
		return
	}

	if os.Args[1] == "test" {
		result := testFunc()
		exitFunc(result)
		return
	} else if os.Args[1] == "install" {
		result := installFunc()
		exitFunc(result)
		return
	} else if os.Args[1] == "uninstall" {
		result := uninstallFunc()
		exitFunc(result)
		return
	}

	commitMessageFile := os.Args[1]
	if commitMessageFile == "" {
		fmt.Print(errors.New("no commit message file passed as parameter 1"))
		exitFunc(1)
		return
	}

	projectConfiguration, err := config.LoadProjectConfigurationFromCommitMessageFileDir(commitMessageFile)
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
