package config

import (
	"path/filepath"

	"github.com/Oppodelldog/filediscovery"
)

const configFilename = "git-commit-hook.yaml"

//LoadConfiguration loads the git-commit-hook configuration from file.
func LoadConfiguration() (*Configuration, error) {
	filePath, err := FindConfigurationFilePath()
	if err != nil {
		return nil, err
	}

	return parse(filePath)
}

//FindConfigurationFilePath searches the git-commit-hook config file in various places and returns
// the first found filepath or error
func FindConfigurationFilePath() (string, error) {

	return filediscovery.New([]filediscovery.FileLocationProvider{
		filediscovery.WorkingDirProvider(".git", "hooks"),
		filediscovery.WorkingDirProvider(".git"),
		filediscovery.WorkingDirProvider(),
		filediscovery.ExecutableDirProvider(),
		filediscovery.HomeConfigDirProvider(".config", "git-commit-hook"),
	}).Discover(configFilename)
}

// LoadProjectConfigurationByName loads a project configuration by its name
func LoadProjectConfigurationByName(projectName string) (Project, error) {

	configuration, err := LoadConfiguration()
	if err != nil {
		return Project{}, err
	}

	return configuration.GetProjectByName(projectName)
}

// LoadProjectConfigurationFromCommitMessageFileDir loads a project configuration by resolving a given commit-message file
func LoadProjectConfigurationFromCommitMessageFileDir(commitMessageFile string) (Project, error) {

	projectPath, err := filepath.Abs(filepath.Dir(commitMessageFile))
	if err != nil {
		return Project{}, err
	}

	configuration, err := LoadConfiguration()
	if err != nil {
		return Project{}, err
	}

	return configuration.GetProjectByRepoPath(projectPath)
}
