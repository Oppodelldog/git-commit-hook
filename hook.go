package gitcommithook

import "fmt"

func ModifyGitCommitMessage(gitCommitMessage string) string {
	branchName := getCurrentGitBranchName()
	return fmt.Sprintf("%s: %s", branchName, gitCommitMessage)
}

func getCurrentGitBranchName() string {
	return "BRANCH"
}
