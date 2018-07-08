package hook

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/Oppodelldog/git-commit-hook/config"
	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
)

type (
	// CommitMessageValidator implementa validation for a given commit message
	CommitMessageValidator interface {
		Validate(branchName, commitMessage string) error
	}

	commitMessageValidator struct {
		projectConfig config.Project
	}
)

// NewCommitMessageValidator creates a new CommitMessageValidator
func NewCommitMessageValidator(projectConfig config.Project) CommitMessageValidator {
	return &commitMessageValidator{
		projectConfig: projectConfig,
	}
}

// Validate validates the given commitMessage. It uses the validations configured for the given branchName.
// As soon as one validation check succeeds, the validation passes.
func (v *commitMessageValidator) Validate(branchName, commitMessage string) error {

	branchType := v.projectConfig.GetBranchType(branchName)
	validators := v.projectConfig.GetValidator(branchType)
	if len(validators) == 0 {
		return nil
	}

	for validationPattern := range validators {
		if regexMatchesString(validationPattern, commitMessage) {
			return nil
		}
	}

	return prepareError(branchName, validators)
}

func prepareError(branchName string, validators map[string]string) error {
	buffer := bytes.NewBufferString("validation error for branch ")
	buffer.WriteString(fmt.Sprintf("'%s'\n", branchName))
	buffer.WriteString("at least expected one of the following to match\n")

	for _, validationDescription := range validators {
		buffer.WriteString(" - ")
		buffer.WriteString(validationDescription)
		buffer.WriteString("\n")
	}

	return errors.New(buffer.String())
}

func regexMatchesString(pattern string, branchName string) bool {
	return pcre.MustCompile(pattern, 0).MatcherString(branchName, 0).Matches()
}
