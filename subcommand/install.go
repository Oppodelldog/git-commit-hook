package subcommand

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/Oppodelldog/git-commit-hook/config"
)

// Install subcommand installs the git-commit-hook in configured git repositories
func Install() int {
	var projectName string
	var allFlag bool
	var forceOverwrite bool
	flagSet := flag.NewFlagSet("git-commit-hook install", flag.ContinueOnError)
	flagSet.StringVar(&projectName, "p", "", `project name`)
	flagSet.BoolVar(&allFlag, "a", false, `all`)
	flagSet.BoolVar(&forceOverwrite, "f", false, `force file creation by overwriting`)
	err := flagSet.Parse(os.Args[2:])
	if err != nil {
		fmt.Printf("git-commit-hook install error: %v\n", err)
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
		err = installForProject(projectConfiguraiton.Path, forceOverwrite)
		if err != nil {
			fmt.Println(err)
			return 1
		}
	} else if allFlag {
		err := installForAllProjects(configuration, forceOverwrite)
		if err != nil {
			fmt.Println(err)
			return 1
		}
	} else {
		flagSet.PrintDefaults()
	}

	return 0
}

func installForAllProjects(configuration *config.Configuration, forceOverwrite bool) error {
	var hasErrors bool
	for _, projectConfiguration := range *configuration {
		err := installForProject(projectConfiguration.Path, forceOverwrite)
		if err != nil {
			hasErrors = true
		}
	}
	if hasErrors {
		return errors.New("done with errors")
	}
	return nil
}

func installForProject(gitFolderPath string, forceOverwrite bool) error {

	exeFile, err := os.Executable()
	if err != nil {
		return err
	}

	commitHookFilePath := createCommitHookFilePath(gitFolderPath)

	fmt.Printf("installing git-commit-hook to '%s': ", commitHookFilePath)

	if _, err = os.Stat(commitHookFilePath); err == nil {
		if !forceOverwrite {
			fmt.Println("file already exists, use -f to force overwriting")
			return nil
		}
		err = os.Remove(commitHookFilePath)
		if err != nil {
			fmt.Println(err)
			return nil
		}
	}

	err = os.Symlink(exeFile, commitHookFilePath)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println("OK")
	return nil
}
