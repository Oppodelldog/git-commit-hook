package subcommand

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/Oppodelldog/git-commit-hook/config"
)

func NewDiagCommand() *DiagCommand {
	return &DiagCommand{
		stdoutWriter:                         io.Writer(os.Stdout),
		findConfigurationFilePath:            config.FindConfigurationFilePath,
		loadConfiguration:                    config.LoadConfiguration,
		checkIsCommitHookInstalledAtPath:     isCommitHookInstalled,
		checkIsAnotherGitHookInstalledAtPath: isAnotherGitHookInstalled,
	}
}

type DiagCommand struct {
	stdoutWriter                         io.Writer
	findConfigurationFilePath            func() (string, error)
	loadConfiguration                    func() (*config.Configuration, error)
	checkIsCommitHookInstalledAtPath     func(string) bool
	checkIsAnotherGitHookInstalledAtPath func(string) bool
}

// Diagnostics gives useful output about the current configuration
func (cmd *DiagCommand) Diagnostics() int {

	configurationFilePath, err := cmd.findConfigurationFilePath()
	if err != nil {
		cmd.stdoutf("error while searching config file: %v\n", err)
		return 1
	}

	cmd.stdout("git-commit-hook diagnostics")
	cmd.stdoutf("load configuration: %v\n", configurationFilePath)
	cmd.stdout("")

	configuration, err := cmd.loadConfiguration()
	if err != nil {
		cmd.stdoutf("error loading configuration: %v\n", err)
		return 1
	}

	for projectName, projectConfiguration := range *configuration {
		cmd.stdout("-------------------------------------------------------------------\n")
		cmd.printProjectConfiguration(projectName, projectConfiguration)

		cmd.stdout("\ngit-commit-hook installed: ")
		if cmd.checkIsCommitHookInstalledAtPath(projectConfiguration.Path) {
			cmd.stdout("YES")
			cmd.stdout("\n")
		} else {
			cmd.stdout("NO")
			if cmd.checkIsAnotherGitHookInstalledAtPath(projectConfiguration.Path) {
				cmd.stdout(", another commit-msg hook is installed")
			}
			cmd.stdout("\n")
		}
	}

	return 0
}

func (cmd *DiagCommand) printProjectConfiguration(projectName string, projectConfiguration config.Project) {
	cmd.stdout("project:", projectName)
	cmd.stdout("path   :", projectConfiguration.Path)
	cmd.stdout("\nbranch types:")
	cmd.printBranchTypes(projectConfiguration.BranchTypes)
	cmd.stdout("\nbranch type templates:")
	cmd.printBranchTemplates(projectConfiguration.Templates)
	cmd.stdout("\nbranch type validation:")
	cmd.printBranchValidation(projectConfiguration.Validation)
}

func (cmd *DiagCommand) printBranchTypes(m map[string]config.BranchTypePattern) {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		cmd.stdout("\t", k, ":", m[k])
	}
}

func (cmd *DiagCommand) printBranchTemplates(m map[string]config.BranchTypeTemplate) {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		cmd.stdout("\t", k, ":", m[k])
	}
}

func (cmd *DiagCommand) printBranchValidation(m map[string]config.BranchValidationConfiguration) {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		cmd.stdout("\t", k, ":")
		var keys2 []string
		for k2 := range m[k] {
			keys2 = append(keys2, k2)
		}
		sort.Strings(keys)
		for _, k2 := range keys2 {
			cmd.stdout("\t\t", k2, ":", m[k][k2])
		}
	}
}

func (cmd *DiagCommand) stdout(i ...interface{}) {
	fmt.Fprint(cmd.stdoutWriter, i...)
}

func (cmd *DiagCommand) stdoutf(format string, i ...interface{}) {
	fmt.Fprintf(cmd.stdoutWriter, format, i...)
}
