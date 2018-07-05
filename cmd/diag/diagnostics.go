package diag

import (
	"fmt"
	"github.com/Oppodelldog/git-commit-hook/config"
	"os"
	"path"
)

func SubCommandDiagnostics() int {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return 1
	}
	projectConfiguration, err := config.LoadProjectConfiguration(path.Join(wd, ".git/commit-msg"))
	if err != nil {
		fmt.Print("no input provided")
		return 1
	}
	fmt.Println("git-commit-hook - parsed configuration")
	fmt.Println("")
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

	return 1
}
