package diag

import (
	"flag"
	"fmt"
	"os"

	"github.com/Oppodelldog/git-commit-hook/config"
	"github.com/Oppodelldog/git-commit-hook/hook"
)

func Test() int {

	var commitMessage string
	var branchName string
	var projectName string

	flagSet := flag.NewFlagSet("test-flagset", flag.ContinueOnError)
	flagSet.StringVar(&commitMessage, "m", "", `commit message`)
	flagSet.StringVar(&branchName, "b", "", `branch name`)
	flagSet.StringVar(&projectName, "p", "", `project name`)
	err := flagSet.Parse(os.Args[2:])
	if err != nil {
		fmt.Printf("git-commit-hook test error: %v\n", err)
		return 1
	}

	if commitMessage != "" {
		configurationFilePath, err := config.FindConfigurationFilePath()
		if err != nil {
			fmt.Printf("error while searching config file: %v\n", err)
			return 1
		}
		var projectConfiguration config.ProjectConfiguration
		if projectName != "" {
			projectConfiguration, err = config.LoadProjectConfigurationByName(projectName)
		} else {
			projectConfiguration, err = loadProjectConfiguration()
		}
		if err != nil {
			fmt.Print(err)
			return 1
		}

		fmt.Printf("testing configuration '%s':\n", configurationFilePath)
		if projectName != "" {
			fmt.Printf("project        : %s\n", projectName)
		}
		if branchName != "" {
			fmt.Printf("branch         : %s\n", branchName)
		}
		fmt.Printf("commit message : %s\n", branchName)
		fmt.Println()

		var modifiedCommitMessage string
		if branchName != "" {
			modifiedCommitMessage, err = hook.ModifyGitCommitMessageForCustomBranch(commitMessage, projectConfiguration, func() (string, error) {
				return branchName, nil
			})
		} else {
			modifiedCommitMessage, err = hook.ModifyGitCommitMessage(commitMessage, projectConfiguration)
		}
		if err != nil {
			fmt.Print(err)
			return 1
		}

		fmt.Println("would generate the following commit message:")
		fmt.Println(modifiedCommitMessage)

	} else {
		flagSet.PrintDefaults()
		return 1
	}

	return 0
}
