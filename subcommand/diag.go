package subcommand

import (
	"sort"

	"os"

	"reflect"

	"github.com/Oppodelldog/git-commit-hook/config"
)

// NewDiagCommand create a new Diag SubCommand
func NewDiagCommand() *DiagCommand {
	return &DiagCommand{
		logger: logger{os.Stdout},
		findConfigurationFilePath:            config.FindConfigurationFilePath,
		loadConfiguration:                    config.LoadConfiguration,
		checkIsCommitHookInstalledAtPath:     isCommitHookInstalled,
		checkIsAnotherGitHookInstalledAtPath: isAnotherGitHookInstalled,
	}
}

// DiagCommand holds data and implementation of the 'diag' sub command
type DiagCommand struct {
	logger
	findConfigurationFilePath            func() (string, error)
	loadConfiguration                    func() (*config.Configuration, error)
	checkIsCommitHookInstalledAtPath     func(string) bool
	checkIsAnotherGitHookInstalledAtPath func(string) bool
}

// Diagnostics gives useful output about the current configuration
func (cmd *DiagCommand) Diagnostics() int {

	configurationFilePath, err := cmd.findConfigurationFilePath()
	if err != nil {
		cmd.stdoutf("error while searching configuration file: %v\n", err)
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
	cmd.stdout("\nbranch types:\n")
	cmd.printConfigurationMap(projectConfiguration.BranchTypes)
	cmd.stdout("\nbranch type templates:\n")
	cmd.printConfigurationMap(projectConfiguration.Templates)
	cmd.stdout("\nbranch type validation:\n")
	cmd.printConfigurationMap(projectConfiguration.Validation)
}

func (cmd *DiagCommand) printConfigurationMap(m interface{}) {
	var keys []string

	for _, value := range reflect.ValueOf(m).MapKeys() {
		keys = append(keys, value.String())
	}

	sort.Strings(keys)

	for _, k := range keys {
		switch v := m.(type) {
		case map[string]config.BranchTypePattern:
			cmd.stdout("\t", k, ":", v[k], "\n")
		case map[string]config.BranchTypeTemplate:
			cmd.stdout("\t", k, ":", v[k], "\n")
		case map[string]config.BranchValidationConfiguration:
			cmd.stdout("\t", k, ":", "\n")
			var keys2 []string
			for k2 := range v[k] {
				keys2 = append(keys2, k2)
			}
			sort.Strings(keys)
			for _, k2 := range keys2 {
				cmd.stdout("\t\t", k2, ":", v[k][k2], "\n")
			}
		}
	}
}
