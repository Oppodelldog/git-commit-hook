package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfiguration_GetProjectConfiguration(t *testing.T) {
	cfg := createTestConfiguration()
	path := "/home/project123"
	projectCfg, err := (&cfg).GetProjectByRepoPath(path)
	if err != nil {
		t.Fatalf("Did not expect GetProjectByRepoPath to return an error, but got: %v ", err)
	}

	assert.Exactly(t, cfg["test2"], projectCfg)
}

func TestConfiguration_GetProjectConfigurationByName(t *testing.T) {
	cfg := createTestConfiguration()
	projectName := "test2"
	projectCfg, err := (&cfg).GetProjectByName(projectName)
	if err != nil {
		t.Fatalf("Did not expect GetProjectByName to return an error, but got: %v ", err)
	}

	assert.Exactly(t, cfg[projectName], projectCfg)
}

func TestGetProjectConfigurationByRepoPath_NotFound(t *testing.T) {
	cfg := createTestConfiguration()
	path := "/home/does-not-exist"
	_, err := (&cfg).GetProjectByRepoPath(path)

	assert.Contains(t, err.Error(), "not found")
}

func TestGetProjectConfigurationByName_NotFound(t *testing.T) {
	cfg := createTestConfiguration()
	projectName := "does-not-exist"
	_, err := (&cfg).GetProjectByName(projectName)

	assert.Contains(t, err.Error(), "not found")
}

func createTestConfiguration() Configuration {
	cfg := Configuration{
		"test1": Project{Path: "/home"},
		"test2": Project{Path: "/home/project123"},
		"test3": Project{Path: ""},
	}
	return cfg
}
