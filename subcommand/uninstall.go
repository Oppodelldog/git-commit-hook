package subcommand

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/Oppodelldog/git-commit-hook/config"
)

// Uninstall subcommand uninstalls the git-commit-hook from configured git repositories
func Uninstall() int {
	var projectName string
	var allFlag bool
	flagSet := flag.NewFlagSet("git-commit-hook uninstall", flag.ContinueOnError)
	flagSet.StringVar(&projectName, "p", "", `project name`)
	flagSet.BoolVar(&allFlag, "a", false, `all`)
	err := flagSet.Parse(os.Args[2:])
	if err != nil {
		fmt.Printf("git-commit-hook uninstall error: %v\n", err)
		return 1
	}

	configuration, err := config.LoadConfiguration()
	if err != nil {
		fmt.Println(err)
		return 1
	}

	if projectName != "" {
		projectConfiguraiton, err := configuration.GetProjectByName(projectName)
		if err != nil {
			fmt.Println(err)
			return 1
		}
		err = uninstallForProject(projectConfiguraiton.Path)
		if err != nil {
			fmt.Println(err)
			return 1
		}
	} else if allFlag {
		err := uninstallForAllProject(configuration)
		if err != nil {
			fmt.Println(err)
			return 1
		}
	} else {
		flagSet.PrintDefaults()
	}

	return 0
}

func uninstallForAllProject(configuration *config.Configuration) error {
	var hasErrors bool
	for _, projectConfiguration := range *configuration {
		err := uninstallForProject(projectConfiguration.Path)
		if err != nil {
			hasErrors = true
		}
	}
	if hasErrors {
		return errors.New("done with errors")
	}
	return nil
}

func uninstallForProject(gitFolderPath string) error {

	commitHookFilePath := createCommitHookFilePath(gitFolderPath)

	fmt.Printf("uninstalling git-commit-hook from '%s': ", commitHookFilePath)

	err := os.Remove(commitHookFilePath)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	fmt.Println("OK")
	return nil
}
