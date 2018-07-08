package config

import (
	"github.com/Oppodelldog/git-commit-hook/regexadapter"
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
		if regexadapter.RegexMatchesString(string(branchTypePattern), branchName) {
			return branchType
		}
	}

	return ""
}

// GetValidator returns the validator that matches the given branch type
// if no
func (projConf *Project) GetValidator(branchType string) map[string]string {
	var foundValidators map[string]string
	for configBranchName, validators := range projConf.Validation {
		if configBranchName == branchType || configBranchName == "*" && foundValidators == nil {
			foundValidators = validators
		}
	}

	return foundValidators
}
