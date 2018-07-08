package subcommand

import (
	"flag"
	"os"

	"github.com/Oppodelldog/git-commit-hook/config"
	"github.com/Oppodelldog/git-commit-hook/git"
	"github.com/Oppodelldog/git-commit-hook/hook"
)

func NewTestCommand() *TestCommand {
	return &TestCommand{
		logger: logger{os.Stdout},
		findConfigurationFilePath:              config.FindConfigurationFilePath,
		loadConfiguration:                      config.LoadConfiguration,
		loadProjectConfigurationByName:         config.LoadProjectConfigurationByName,
		loadProjectConfigurationFromWorkingDir: loadProjectConfiguration,
		newCommitMessageModifier:               newCommitMessageModifier,
	}
}

type TestCommand struct {
	logger
	findConfigurationFilePath              func() (string, error)
	loadConfiguration                      func() (*config.Configuration, error)
	loadProjectConfigurationByName         func(string) (config.Project, error)
	loadProjectConfigurationFromWorkingDir func() (config.Project, error)
	newCommitMessageModifier               func(projectConfiguration config.Project) hook.CommitMessageModifier
}

// Test helps to test configuration against manual input to simulate real commit situations
func (cmd *TestCommand) Test() int {

	var commitMessage string
	var branchName string
	var projectName string

	flagSet := flag.NewFlagSet("git-commit-hook test", flag.ContinueOnError)
	flagSet.SetOutput(cmd.stdoutWriter)
	flagSet.StringVar(&commitMessage, "m", "", `commit message`)
	flagSet.StringVar(&branchName, "b", "", `branch name`)
	flagSet.StringVar(&projectName, "p", "", `project name`)
	err := flagSet.Parse(os.Args[2:])
	if err != nil {
		return 1
	}

	configurationFilePath, err := cmd.findConfigurationFilePath()
	if err != nil {
		cmd.stdoutf("error while searching configuration file: %v\n", err)
		return 1
	}

	if commitMessage == "" {
		cmd.stdout("you must at least enter a commit message using parameter -m\n")
		flagSet.Usage()
		return 1
	}

	if branchName == "" {
		branchName, err = git.GetCurrentBranchName()
		if err != nil {
			cmd.stdout("error while reading branch name. ensure working dir is a git repo or use parameter -b to simulate a branch name\n")
			return 1
		}

		branchName += " (current git branch)"
	}

	var projectConfiguration config.Project
	if projectName != "" {
		projectConfiguration, err = cmd.loadProjectConfigurationByName(projectName)
	} else {
		projectConfiguration, err = cmd.loadProjectConfigurationFromWorkingDir()
	}
	if err != nil {
		cmd.stdout(err, "\n")
		return 1
	}

	cmd.stdoutf("testing configuration '%s':\n", configurationFilePath)
	if projectName != "" {
		cmd.stdoutf("project        : %s\n", projectName)
	}
	cmd.stdoutf("branch         : %s\n", branchName)
	cmd.stdoutf("commit message : %s\n", commitMessage)
	cmd.stdout("\n")

	var modifiedCommitMessage string

	commitMessageModifier := cmd.newCommitMessageModifier(projectConfiguration)
	modifiedCommitMessage, err = commitMessageModifier.ModifyGitCommitMessage(commitMessage, branchName)
	if err != nil {
		cmd.stdout(err, "\n")
		return 1
	}

	cmd.stdoutf("would generate the following commit message:\n%v\n", modifiedCommitMessage)

	return 0
}

func newCommitMessageModifier(projectConfiguration config.Project) hook.CommitMessageModifier {
	return hook.NewCommitMessageModifier(projectConfiguration)
}
