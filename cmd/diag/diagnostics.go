package diag

import (
	"fmt"

	"github.com/Oppodelldog/git-commit-hook/config"
)

func SubCommandDiagnostics() int {

	configurationFilePath, err := config.FindConfigurationFilePath()
	if err != nil {
		fmt.Printf("error while searching config file: %v\n", err)
		return 1
	}

	fmt.Println("git-commit-hook diagnostics")
	fmt.Printf("load configuration: %v\n", configurationFilePath)
	fmt.Println("")

	configuration, err := config.LoadConfiguration()
	if err != nil {
		fmt.Printf("error loading configuration: %v\n", err)
		return 1
	}

	for projectName, projectConfiguration := range *configuration {
		fmt.Println("-------------------------------------------------------------------")
		printProjectConfiguration(projectName, projectConfiguration)

		fmt.Print("\ngit-commit-hook installed: ")
		if isCommitHookInstalled(projectConfiguration.Path) {
			fmt.Println("YES")
		} else {
			fmt.Print("NO")
			if isAnotherGitHookInstalled(projectConfiguration.Path) {
				fmt.Print(", another commit-msg hook is installed")
			}
			fmt.Println()
		}
	}

	return 1
}

func printProjectConfiguration(projectName string, projectConfiguration config.ProjectConfiguration) {
	fmt.Println("project:", projectName)
	fmt.Println("path   :", projectConfiguration.Path)
	fmt.Println("\nbranch types:")
	for branchType, pattern := range projectConfiguration.BranchTypes {
		fmt.Println("\t", branchType, ":", pattern)
	}
	fmt.Println("\nbranch type templates:")
	for branchType, template := range projectConfiguration.Templates {
		fmt.Println("\t", branchType, ":", template)
	}
	fmt.Println("\nbranch type validation:")
	for branchType, validation := range projectConfiguration.Validation {
		fmt.Println("\t", branchType, ":")
		for pattern, description := range validation {
			fmt.Println("\t\t", pattern, ":", description)
		}
	}
}
