package gitcommithook

import (
	"github.com/Oppodelldog/git-commit-hook/git"
	"github.com/pkg/errors"
	"github.com/Oppodelldog/git-commit-hook/config"
	"github.com/Oppodelldog/filediscovery"
)

type gitBranchNameReaderFuncDef func() (string, error)
type featureBranchDetectFuncDef func(string) bool

var featureBranchDetectFunc = featureBranchDetectFuncDef(IsFeatureBranch)
var gitBranchNameReaderFunc = gitBranchNameReaderFuncDef(git.GetCurrentBranchName)

//ModifyGitCommitMessage prepends the current branch name to the given git commit message.
// if the current branch name is detected to be NO feature branch, the user will be prompted to enter
// a feature branch manually. This is then inserted in between current branch and commit message.
// If no valid branch name could be determined the function returns an error
func ModifyGitCommitMessage(gitCommitMessage string) (modifiedCommitMessage string, err error) {

	if gitCommitMessage == "" {
		err = errors.New("commit message is empty")
		return
	}

	branchName, err := gitBranchNameReaderFunc()
	if err != nil {
		return
	}

	if branchName == "" {
		err = errors.New("branch name is empty")
		return
	}

	cfg, err := loadConfiguration()
	if err != nil {
		return
	}

	viewModel := map[string]string{"CommitMessage": gitCommitMessage}
	modifiedCommitMessage, err := cfg.RenderCommitMessage(branchName, viewModel)
	if err != nil {
		return
	}

	err = cfg.Validate(branchName, modifiedCommitMessage)

	return
}

func loadConfiguration() (*config.Configuration, error) {
	const commitHookConfig = "git-commit-hook.yaml"

	filePath, err := filediscovery.New([]filediscovery.FileLocationProvider{
		filediscovery.WorkingDirProvider(),
		filediscovery.ExecutableDirProvider(),
		filediscovery.HomeConfigDirProvider(".config", "git-commit-hook"),
	}).Discover(commitHookConfig)

	if err != nil {
		return nil, err
	}

	return config.Parse(filePath)
}
