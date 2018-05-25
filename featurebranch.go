package gitcommithook

import "regexp"

//noinspection SpellCheckingInspection
var falsePatterns = []string{
	"^release$",
	"^release/.*$",
	"^master$",
	"^develop",
	"^hotfix/.*$",
}

//IsFeatureBranch returns true if the given branchname matches a valid feature branch
func IsFeatureBranch(branchName string) bool {
	for _, falsePattern := range falsePatterns {
		if matches(falsePattern, branchName) {
			return false
		}
	}

	return true
}

func matches(regEx, input string) bool {
	return regexp.MustCompile(regEx).MatchString(input)
}
