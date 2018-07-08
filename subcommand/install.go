package subcommand

import (
	"errors"
	"flag"
	"os"

	"github.com/Oppodelldog/git-commit-hook/config"
)

// NewInstallCommand creates a new InstallCommand
func NewInstallCommand() *InstallCommand {
	return &InstallCommand{
		logger:            logger{os.Stdout},
		loadConfiguration: config.LoadConfiguration,
		gitHookInstaller:  NewGitHookInstaller(),
	}
}

// InstallCommand holds data and implementation of the install sub command
type InstallCommand struct {
	logger
	loadConfiguration func() (*config.Configuration, error)
	gitHookInstaller  GitHookInstaller
}

// Install subcommand installs the git-commit-hook in configured git repositories
func (cmd *InstallCommand) Install() int {
	var projectName string
	var allFlag bool
	var forceOverwrite bool
	flagSet := flag.NewFlagSet("git-commit-hook install", flag.ContinueOnError)
	flagSet.SetOutput(cmd.stdoutWriter)
	flagSet.StringVar(&projectName, "p", "", `project name`)
	flagSet.BoolVar(&allFlag, "a", false, `all`)
	flagSet.BoolVar(&forceOverwrite, "f", false, `force file creation by overwriting`)
	err := flagSet.Parse(os.Args[2:])
	if err != nil {
		return 1
	}

	configuration, err := cmd.loadConfiguration()
	if err != nil {
		cmd.stdout(err, "\n")
		return 1
	}

	if projectName != "" {
		projectConfiguration, err := configuration.GetProjectByName(projectName)
		if err != nil {
			cmd.stdout(err, "\n")
			return 1
		}
		cmd.stdoutf("installing git-commit-hook to '%s': ", projectConfiguration.Path)
		err = cmd.gitHookInstaller.installForProject(projectConfiguration.Path, forceOverwrite)
		if err != nil {
			cmd.stdout(err, "\n")
			return 1
		}
		cmd.stdout("OK", "\n")

	} else if allFlag {
		err := cmd.installForAllProjects(configuration, forceOverwrite)
		if err != nil {
			cmd.stdout(err, "\n")
			return 1
		}
	} else {
		flagSet.Usage()
		return 1
	}

	return 0
}

func (cmd *InstallCommand) installForAllProjects(configuration *config.Configuration, forceOverwrite bool) error {
	var hasErrors bool
	for _, projectConfiguration := range *configuration {
		cmd.stdoutf("installing git-commit-hook to '%s': ", projectConfiguration.Path)
		err := cmd.gitHookInstaller.installForProject(projectConfiguration.Path, forceOverwrite)
		if err != nil {
			cmd.stdout(err, "\n")
			hasErrors = true
		} else {
			cmd.stdout("OK", "\n")
		}
	}
	if hasErrors {
		return errors.New("done with errors")
	}
	return nil
}
