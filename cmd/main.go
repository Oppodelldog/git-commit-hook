package main

import (
	"fmt"

	"os"

	"github.com/Oppodelldog/git-commit-hook"
)

func main() {
	commitMessage, err := gitcommithook.ModifyGitCommitMessage(os.Args[1])
	if err != nil {
		fmt.Printf("error in git hook: %s", err.Error())
		os.Exit(1)
	}

	fmt.Print(commitMessage)
}
