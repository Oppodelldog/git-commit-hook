package subcommand

import (
	"errors"
	"flag"
	"os"

	"github.com/Oppodelldog/git-commit-hook/config"
)

// UninstallCommand holds data and methods for 'uninstall' subcommand
type UninstallCommand struct {
	logger
	loadConfiguration func() (*config.Configuration, error)
	deleteFile        func(string) error
}

//NewUninstallerCommand creates a new uninstall subcommand
func NewUninstallerCommand() *UninstallCommand {
	return &UninstallCommand{
		logger:            logger{stdoutWriter: os.Stdout},
		loadConfiguration: config.LoadConfiguration,
		deleteFile:        os.Remove,
	}
}

// Uninstall subcommand uninstalls the git-commit-hook from configured git repositories
func (u *UninstallCommand) Uninstall() int {
	var projectName string
	var allFlag bool
	flagSet := flag.NewFlagSet("git-commit-hook uninstall", flag.ContinueOnError)
	flagSet.SetOutput(u.stdoutWriter)
	flagSet.StringVar(&projectName, "p", "", `project name`)
	flagSet.BoolVar(&allFlag, "a", false, `all`)
	err := flagSet.Parse(os.Args[2:])
	if err != nil {
		return 1
	}

	configuration, err := u.loadConfiguration()
	if err != nil {
		u.stdout(err, "\n")
		return 1
	}

	if projectName != "" {
		projectConfiguraiton, err := configuration.GetProjectByName(projectName)
		if err != nil {
			u.stdout(err, "\n")
			return 1
		}
		err = u.uninstallForProject(projectConfiguraiton.Path)
		if err != nil {
			u.stdout(err, "\n")
			return 1
		}
	} else if allFlag {
		err := u.uninstallForAllProject(configuration)
		if err != nil {
			u.stdout(err, "\n")
			return 1
		}
	} else {
		flagSet.Usage()
	}

	return 0
}

func (u *UninstallCommand) uninstallForAllProject(configuration *config.Configuration) error {
	var hasErrors bool
	for _, projectConfiguration := range *configuration {
		err := u.uninstallForProject(projectConfiguration.Path)
		if err != nil {
			hasErrors = true
		}
	}
	if hasErrors {
		return errors.New("done with errors")
	}
	return nil
}

func (u *UninstallCommand) uninstallForProject(gitFolderPath string) error {

	commitHookFilePath := createCommitHookFilePath(gitFolderPath)

	u.stdoutf("uninstalling git-commit-hook from '%s': ", commitHookFilePath)

	err := u.deleteFile(commitHookFilePath)
	if err != nil {
		u.stdout(err, "\n")
		return err
	}

	u.stdout("OK\n")
	return nil
}
