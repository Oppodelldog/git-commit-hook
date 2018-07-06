package config

import (
	"path/filepath"

	"github.com/Oppodelldog/filediscovery"
)

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
	const configFilename = "git-commit-hook.yaml"

	return filediscovery.New([]filediscovery.FileLocationProvider{
		filediscovery.WorkingDirProvider(),
		filediscovery.WorkingDirProvider(".git"),
		filediscovery.WorkingDirProvider(".git", "hooks"),
		filediscovery.ExecutableDirProvider(),
		filediscovery.HomeConfigDirProvider(".config", "git-commit-hook"),
	}).Discover(configFilename)
}

func LoadProjectConfigurationByName(projectName string) (ProjectConfiguration, error) {

	configuration, err := LoadConfiguration()
	if err != nil {
		return ProjectConfiguration{}, err
	}

	return configuration.GetProjectConfigurationByName(projectName)
}

func LoadProjectConfigurationFromCommitMessageFileDir(commitMessageFile string) (ProjectConfiguration, error) {

	projectPath, err := filepath.Abs(filepath.Dir(commitMessageFile))
	if err != nil {
		return ProjectConfiguration{}, err
	}

	configuration, err := LoadConfiguration()
	if err != nil {
		return ProjectConfiguration{}, err
	}

	return configuration.GetProjectConfiguration(projectPath)
}
