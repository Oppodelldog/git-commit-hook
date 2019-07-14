package subcommand

import (
	"errors"
	"os"
)

// NewGitHookInstaller creates a new GitHookInstaller
func NewGitHookInstaller() GitHookInstaller {
	return &gitHookInstaller{
		logger:                       logger{stdoutWriter: os.Stdout},
		getCurrentExecutableFilePath: os.Executable,
		existsFile:                   checkFileExists,
		removeFile:                   os.Remove,
		createSymlink:                os.Symlink,
	}
}

type (
	// GitHookInstaller implements the installation of git-commit-hook into a repository
	GitHookInstaller interface {
		installForProject(gitFolderPath string, forceOverwrite bool) error
	}
	gitHookInstaller struct {
		logger
		getCurrentExecutableFilePath func() (string, error)
		existsFile                   func(string) bool
		removeFile                   func(string) error
		createSymlink                func(string, string) error
	}
)

func (cmd *gitHookInstaller) installForProject(gitFolderPath string, forceOverwrite bool) error {

	exeFile, err := cmd.getCurrentExecutableFilePath()
	if err != nil {
		return err
	}

	commitHookFilePath := createCommitHookFilePath(gitFolderPath)

	if cmd.existsFile(commitHookFilePath) {
		if !forceOverwrite {
			return errors.New("file already exists, use -f to force overwriting")
		}

		err = cmd.removeFile(commitHookFilePath)
		if err != nil {
			return err
		}
	}

	return cmd.createSymlink(exeFile, commitHookFilePath)
}

func checkFileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}
