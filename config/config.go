package config

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"

	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
)

type (
	Configuration map[string]ProjectConfiguration

	ProjectConfiguration struct {
		Path        string                                   `yaml:"path"`
		BranchTypes map[string]BranchTypeConfiguration       `yaml:"branch"`
		Templates   map[string]BranchTemplateConfiguration   `yaml:"template"`
		Validation  map[string]BranchValidationConfiguration `yaml:"validation"`
	}

	BranchValidationConfiguration map[string]string

	BranchTemplateConfiguration struct {
		Template string `yaml:"template"`
	}

	BranchTypeConfiguration struct {
		Pattern string `yaml:"matcher"`
	}
)

func (configuration *Configuration) GetProjectConfiguration(path string) (ProjectConfiguration, error) {
	for _, projectCfg := range *configuration {
		if projectCfg.Path == path {
			return projectCfg, nil
		}
	}

	return ProjectConfiguration{}, fmt.Errorf("project configuration not found for path %s", path)
}

type ViewModel struct {
	BranchName    string
	CommitMessage string
}

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

func (projConf *ProjectConfiguration) RenderCommitMessage(branchName string, viewModel ViewModel) (string, error) {
	branchType := projConf.GetBranchType(branchName)
	commitMessageTemplate := projConf.GetTemplate(branchType)
	tmpl, err := template.New("commitMessageTemplate").Parse(commitMessageTemplate)
	if err != nil {
		return "", err
	}
	buffer := bytes.NewBufferString("")
	err = tmpl.Execute(buffer, viewModel)

	return buffer.String(), err
}

func (projConf *ProjectConfiguration) GetTemplate(branchType string) string {
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

func (projConf *ProjectConfiguration) Validate(branchName string, commitMessage string) error {

	branchType := projConf.GetBranchType(branchName)
	validators := projConf.getValidators(branchType)
	if validators == nil || len(validators) == 0 {
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
