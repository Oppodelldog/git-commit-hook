package config

import (
	"text/template"
	"bytes"
	"fmt"
)

type
(
	Configuration map[string]ProjectConfiguration

	ProjectConfiguration struct {
		Path       string                                   `yaml:"path"`
		Branches   map[string]BranchMatcherConfiguration    `yaml:"branch"`
		Templates  map[string]BranchTemplateConfiguration   `yaml:"template"`
		Validation map[string]BranchValidationConfiguration `yaml:"validation"`
	}

	BranchValidationConfiguration map[string]string

	BranchTemplateConfiguration struct {
		Template string `yaml:"template"`
	}

	BranchMatcherConfiguration struct {
		Matcher string `yaml:"matcher"`
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

func (projConf *ProjectConfiguration) RenderCommitMessage(branchName string, viewModel ViewModel) (string, error) {
	commitMessageTemplate := projConf.GetTemplate(branchName)
	tmpl, err := template.New("commitMessageTemplate").Parse(commitMessageTemplate)
	if err != nil {
		return "", err
	}
	buffer := bytes.NewBufferString("")
	err = tmpl.Execute(buffer, viewModel)

	return buffer.String(), err
}

func (projConf *ProjectConfiguration) GetTemplate(branchName string) string {
	foundTemplate := ""
	for configBranchName, branchTemplateCfg := range projConf.Templates {
		if configBranchName == branchName || configBranchName == "*" && foundTemplate == "" {
			foundTemplate = branchTemplateCfg.Template
		}
	}

	return foundTemplate
}
