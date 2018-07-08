package config

import "fmt"

type (
	// Configuration is the data representation of the config file structure
	Configuration map[string]Project

	// BranchTypeTemplate holds a go template that defined the modified commit message
	BranchTypeTemplate string

	// BranchTypePattern defines a regex pattern that identifies a branch type
	BranchTypePattern string

	// BranchValidationConfiguration holds a regex pattern in the index and a test description in the value of the map
	// The pattern will be tested against a prepared commit message.
	BranchValidationConfiguration map[string]string
)

// GetProjectByRepoPath returns a Project for the given git repository path
func (c *Configuration) GetProjectByRepoPath(path string) (Project, error) {
	for _, projectCfg := range *c {
		if projectCfg.Path == path {
			return projectCfg, nil
		}
	}

	return Project{}, fmt.Errorf("project configuration not found for path '%s'", path)
}

// GetProjectByName returns a Project for the given project name
func (c *Configuration) GetProjectByName(projectName string) (Project, error) {
	for configProjectName, projectCfg := range *c {
		if configProjectName == projectName {
			return projectCfg, nil
		}
	}

	return Project{}, fmt.Errorf("project configuration not found for project name '%s'", projectName)
}
