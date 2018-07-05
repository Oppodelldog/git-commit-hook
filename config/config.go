package config

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"

	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
)

type (
	// Configuration is the data representation of the config file structure
	Configuration map[string]ProjectConfiguration

	// ProjectConfiguration defined a project related section of the Configuration
	ProjectConfiguration struct {
		// Path to the git repository this configuration should be used while committing
		Path string `yaml:"path"`
		// BranchTypes is a map whose key defines a BranchType - it's value holds a pattern that identifies
		// a given branch name to be of that branch type.
		BranchTypes map[string]BranchTypeConfiguration `yaml:"branch"`
		// Templates is a map whose key refers a branchType - it's value holds a go template that will render the commit message
		Templates map[string]BranchTemplateConfiguration `yaml:"template"`
		// Validation is a map whose key refers a branchType - it's value holds configuration to validate the created commit message
		Validation map[string]BranchValidationConfiguration `yaml:"validation"`
	}

	// BranchTemplateConfiguration holds a go template that defined the modified commit message
	BranchTemplateConfiguration struct {
		Template string `yaml:"template"`
	}

	// BranchTypeConfiguration holds a pattern that identifies a branch type
	BranchTypeConfiguration struct {
		Pattern string `yaml:"matcher"`
	}

	// BranchValidationConfiguration holds a regex pattern in the index and a test description in the value of the map
	// The pattern will be tested against a prepared commit message.
	BranchValidationConfiguration map[string]string
)

// GetProjectConfiguration returns a ProjectConfiguration for the given git repository path
func (configuration *Configuration) GetProjectConfiguration(path string) (ProjectConfiguration, error) {
	for _, projectCfg := range *configuration {
		if projectCfg.Path == path {
			return projectCfg, nil
		}
	}

	return ProjectConfiguration{}, fmt.Errorf("project configuration not found for path %s", path)
}

// ViewModel defines all variables that can be in templates to define the modified commit message
type ViewModel struct {
	BranchName    string
	CommitMessage string
}

// GetBranchType returns a branch type for the given branch name or empty string if no branch type was found.
func (projConf *ProjectConfiguration) GetBranchType(branchName string) string {

	for branchType, matcher := range projConf.BranchTypes {
		if regexMatchesString(matcher.Pattern, branchName) {
			return branchType
		}
	}

	return ""
}

func regexMatchesString(pattern string, branchName string) bool {
	return pcre.MustCompile(pattern, 0).MatcherString(branchName, 0).Matches()
}

// RenderCommitMessage renders a commit message using the template defined for the given branchName
func (projConf *ProjectConfiguration) RenderCommitMessage(branchName string, viewModel ViewModel) (string, error) {
	branchType := projConf.GetBranchType(branchName)
	commitMessageTemplate := projConf.getTemplate(branchType)
	tmpl, err := template.New("commitMessageTemplate").Parse(commitMessageTemplate)
	if err != nil {
		return "", err
	}
	buffer := bytes.NewBufferString("")
	err = tmpl.Execute(buffer, viewModel)

	return buffer.String(), err
}

func (projConf *ProjectConfiguration) getTemplate(branchType string) string {
	foundTemplate := ""
	for configBranchName, branchTemplateCfg := range projConf.Templates {
		if configBranchName == branchType || configBranchName == "*" && foundTemplate == "" {
			foundTemplate = branchTemplateCfg.Template
		}
	}

	return foundTemplate
}

func (projConf *ProjectConfiguration) getValidators(branchType string) map[string]string {
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
func (projConf *ProjectConfiguration) Validate(branchName string, commitMessage string) error {

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
		buffer.WriteString(validationDescription)
		buffer.WriteString("\n")
	}

	return errors.New(buffer.String())
}
