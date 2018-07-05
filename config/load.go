package config

import (
	"github.com/Oppodelldog/filediscovery"
	"path/filepath"
)

//LoadConfiguration loads the git-commit-hook configuration from file.
func LoadConfiguration() (*Configuration, error) {
	const commitHookConfig = "git-commit-hook.yaml"

	filePath, err := filediscovery.New([]filediscovery.FileLocationProvider{
		filediscovery.WorkingDirProvider(),
		filediscovery.WorkingDirProvider(".git"),
		filediscovery.WorkingDirProvider(".git", "hooks"),
		filediscovery.ExecutableDirProvider(),
		filediscovery.HomeConfigDirProvider(".config", "git-commit-hook"),
	}).Discover(commitHookConfig)

	if err != nil {
		return nil, err
	}

	return parse(filePath)
}

func LoadProjectConfiguration(commitMessageFile string) (ProjectConfiguration, error) {

	configuration, err := LoadConfiguration()
	if err != nil {
		return ProjectConfiguration{}, err
	}

	projectPath, err := filepath.Abs(filepath.Dir(commitMessageFile))
	if err != nil {
		return ProjectConfiguration{}, err
	}

	return configuration.GetProjectConfiguration(projectPath)
}
