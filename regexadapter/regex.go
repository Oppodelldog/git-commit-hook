package regexadapter

import (
	"regexp"
)

func RegexMatchesString(pattern string, branchName string) bool {
	return regexp.MustCompile(pattern).MatchString(branchName)
	//return pcre.MustCompile(pattern, 0).MatcherString(branchName, 0).Matches()
}

