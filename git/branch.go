package git

import (
	"os/exec"
	"regexp"
)

type execFuncDef func(s1 string, s2 ...string) *exec.Cmd

var execFunc = execFuncDef(exec.Command)

//GetCurrentBranchName executes 'git branch' to get the current branch
func GetCurrentBranchName() (string, error) {
	outputBytes, err := execFunc("git", "branch").Output()
	if err != nil {
		return "", err
	}

	branchName := getBranchNameFromGitOutput(string(outputBytes))

	return branchName, err
}

func getBranchNameFromGitOutput(gitOutput string) string {

	var re = regexp.MustCompile(`(?m)^\* (.*)$`)

	matches := re.FindAllStringSubmatch(gitOutput, 1)

	if len(matches) > 0 {
		if len(matches[0]) > 1 {
			return matches[0][1]
		}
	}

	return ""
}
