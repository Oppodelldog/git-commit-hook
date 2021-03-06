package subcommand

import (
	"fmt"
	"os"
	"path"

	"path/filepath"

	"github.com/Oppodelldog/git-commit-hook/config"
)

const gitCommitMessageHookName = "commit-msg"

func loadProjectConfiguration() (config.Project, error) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return config.Project{}, err
	}

	return config.LoadProjectConfigurationFromCommitMessageFileDir(path.Join(wd, ".git/commit-message.txt"))
}

func createCommitHookFilePath(gitFolderPath string) string {
	commitHookFilePath := path.Join(gitFolderPath, "hooks", gitCommitMessageHookName)

	return commitHookFilePath
}

func isAnotherGitHookInstalled(gitFolderPath string) bool {
	commitHookFilePath := createCommitHookFilePath(gitFolderPath)
	if _, err := os.Stat(commitHookFilePath); err == nil {
		return true
	}

	return false
}

func isCommitHookInstalled(gitFolderPath string) bool {
	commitHookFilePath := createCommitHookFilePath(gitFolderPath)
	commitHookOrignFilePath, err := filepath.EvalSymlinks(commitHookFilePath)
	if err != nil {
		return false
	}

	exeFilePath, err := os.Executable()
	if err != nil {
		return false
	}

	if exeFilePath == commitHookOrignFilePath {
		return true
	}

	return false
}
