package config

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
)

type (
	// Project defined a project related section of the Configuration
	Project struct {
		// Path to the git repository this configuration should be used while committing
		Path string `yaml:"path"`
		// BranchTypes is a map whose key defines a BranchType - it's value holds a pattern that identifies
		// a given branch name to be of that branch type.
		BranchTypes map[string]BranchTypePattern `yaml:"branch"`
		// Templates is a map whose key refers a branchType - it's value holds a go template that will render the commit message
		Templates map[string]BranchTypeTemplate `yaml:"template"`
		// Validation is a map whose key refers a branchType - it's value holds configuration to validate the created commit message
		Validation map[string]BranchValidationConfiguration `yaml:"validation"`
	}
)

// GetBranchType returns a branch type for the given branch name or empty string if no branch type was found.
func (projConf *Project) GetBranchType(branchName string) string {

	for branchType, branchTypePattern := range projConf.BranchTypes {
		if regexMatchesString(string(branchTypePattern), branchName) {
			return branchType
		}
	}

	return ""
}

func regexMatchesString(pattern string, branchName string) bool {
	return pcre.MustCompile(pattern, 0).MatcherString(branchName, 0).Matches()
}

func (projConf *Project) getValidators(branchType string) map[string]string {
	var foundValidators map[string]string
	for configBranchName, validators := range projConf.Validation {
		if configBranchName == branchType || configBranchName == "*" && foundValidators == nil {
			foundValidators = validators
		}
	}

	return foundValidators
}

// Validate validates the given commitMessage. It uses the validations configured for the given branchName.
// As soon as one validation check succeeds, the validation passes.
func (projConf *Project) Validate(branchName string, commitMessage string) error {

	branchType := projConf.GetBranchType(branchName)
	validators := projConf.getValidators(branchType)
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
