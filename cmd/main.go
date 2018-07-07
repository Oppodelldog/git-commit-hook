package main

import (
	"os"

	"fmt"

	"github.com/Oppodelldog/git-commit-hook/config"
	"github.com/Oppodelldog/git-commit-hook/hook"
	"github.com/Oppodelldog/git-commit-hook/subcommand"
	"github.com/pkg/errors"
)

type callWithIntResult func() int
type rewriteCommitMessageFuncDef func(string, hook.CommitMessageModifier) error
type exitFuncDef func(code int)

var diagnosticsFunc = callWithIntResult(subcommand.NewDiagCommand().Diagnostics)
var testFunc = callWithIntResult(subcommand.NewTestCommand().Test)
var installFunc = callWithIntResult(subcommand.Install)
var uninstallFunc = callWithIntResult(subcommand.Uninstall)
var rewriteCommitMessageFunc = rewriteCommitMessageFuncDef(hook.RewriteCommitMessage)
var exitFunc = exitFuncDef(os.Exit)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("too few arguments")
		fmt.Println("")
		fmt.Println("diag		- gives diagnostic output of configuration")
		fmt.Println("install 	- helps to install git-commit-hook")
		fmt.Println("uninstall 	- helps to uninstall git-commit-hook")
		fmt.Println("test 		- helps to test configuration with manual inputs")
		exitFunc(0)
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
	} else if os.Args[1] == "diag" {
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

	projectConfiguration, err := config.LoadProjectConfigurationFromCommitMessageFileDir(commitMessageFile)
	if err != nil {
		fmt.Print(err)
		exitFunc(1)
		return
	}

	commitMessageModifier := hook.NewCommitMessageModifier(projectConfiguration)
	err = rewriteCommitMessageFunc(commitMessageFile, commitMessageModifier)
	if err != nil {
		fmt.Print(err)
		exitFunc(1)
		return
	}

	exitFunc(0)
}
