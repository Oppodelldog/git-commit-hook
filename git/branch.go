package git

import (
	"os/exec"
	"regexp"
)

type execFuncDef func(s1 string, s2 ...string) *exec.Cmd

var execFunc = execFuncDef(exec.Command)

//GetCurrentBranchName executes 'git branch' to get the current branch
func GetCurrentBranchName() (string, error) {
	outputBytes, err := execFunc("git", "branch").CombinedOutput()
	if err != nil {
		return "", err
	}

	branchName := getBranchNameFromGitOutput(string(outputBytes))

	if branchName == "" {
		outputBytes, err := execFunc("git", "log").CombinedOutput()
		if err != nil {
			return "", err
		}

		branchName = getBranchNameFromGitLogOutput(string(outputBytes))
	}

	return branchName, nil
}

func getBranchNameFromGitOutput(gitOutput string) string {

	return extractFromString(gitOutput, `(?m)^\* (.*)$`)
}

func getBranchNameFromGitLogOutput(gitOutput string) string {
	// https://github.com/git/git/blob/ed843436dd4924c10669820cc73daf50f0b4dabd/revision.c#L2303
	pattern := `(?m)^fatal: your current branch '(.*)' does not have any commits yet$`

	return extractFromString(gitOutput, pattern)
}

func extractFromString(s, regexPattern string) string {
	var re = regexp.MustCompile(regexPattern)

	matches := re.FindAllStringSubmatch(s, 1)

	if len(matches) > 0 {
		if len(matches[0]) > 1 {
			return matches[0][1]
		}
	}

	return ""
}
