package main

import (
	"fmt"

	"os"

	"github.com/Oppodelldog/git-commit-hook"
)

type rewriteCommitMessageFuncDef func(string) error
type exitFuncDef func(code int)

var rewriteCommitMessageFunc = rewriteCommitMessageFuncDef(gitcommithook.RewriteCommitMessage)
var exitFunc = exitFuncDef(os.Exit)

func main() {
	commitMessageFile := os.Args[1]
	err := rewriteCommitMessageFunc(commitMessageFile)
	if err != nil {
		fmt.Print(err)
		exitFunc(1)
		return
	}
	exitFunc(0)
}
