package gitcommithook

import "regexp"

var patterns = []string{
	`(?m)(?:\s|^|/)([A-Z]+-[0-9]+)([\s,;:!.-]|$)`,
}

//IsFeatureBranch detects weatcher the input matches one of the defined patterns
func IsFeatureBranch(commitMessage string) bool {
	for _, pattern := range patterns {
		if matches(pattern, commitMessage) {
			return true
		}
	}

	return false
}

func matches(regEx, input string) bool {
	return regexp.MustCompile(regEx).MatchString(input)
}
