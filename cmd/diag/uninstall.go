package diag

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/Oppodelldog/git-commit-hook/config"
)

func Uninstall() int {
	var projectName string
	var allFlag bool
	flagSet := flag.NewFlagSet("uninstall-flagset", flag.ContinueOnError)
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
		projectConfiguraiton, err := configuration.GetProjectConfigurationByName(projectName)
		if err != nil {
			fmt.Println(err)
			return 1
		}
		err = uninstallProject(projectConfiguraiton.Path)
		if err != nil {
			fmt.Println(err)
			return 1
		}
	} else if allFlag {
		err := uninstallAllProject(configuration)
		if err != nil {
			fmt.Println(err)
			return 1
		}
	} else {
		flagSet.PrintDefaults()
	}

	return 0
}

func uninstallAllProject(configuration *config.Configuration) error {
	var hasErrors bool
	for _, projectConfiguration := range *configuration {
		err := uninstallProject(projectConfiguration.Path)
		if err != nil {
			hasErrors = true
		}
	}
	if hasErrors {
		return errors.New("done with errors")
	}
	return nil
}

func uninstallProject(gitFolderPath string) error {

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
