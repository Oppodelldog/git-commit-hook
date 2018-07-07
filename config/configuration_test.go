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

func TestValidate(t *testing.T) {

	testDataSet := []struct {
		BranchName    string
		CommitMessage string
		ErrorContains string
	}{
		{
			BranchName:    "branch1",
			CommitMessage: "commit message that contains A.",
			ErrorContains: "",
		},
		{
			BranchName:    "branch1",
			CommitMessage: "commit message that contains B.",
			ErrorContains: "",
		},
		{
			BranchName:    "branch2",
			CommitMessage: "commit message that contains A.",
			ErrorContains: "",
		},
		{
			BranchName:    "branch2",
			CommitMessage: "commit message that contains B.",
			ErrorContains: "",
		},
		{
			BranchName:    "branch3",
			CommitMessage: "does not matter what this contains, since there are no validator for branch3",
			ErrorContains: "",
		},
		{
			BranchName:    "branch1",
			CommitMessage: "commit message that contains Z.",
			ErrorContains: "validation error for branch 'branch1'",
		},
	}

	for testName, testData := range testDataSet {
		t.Run(string(testName), func(t *testing.T) {

			cfg := &Project{
				BranchTypes: map[string]BranchTypePattern{
					"branch1": `^branch1$`,
					"branch2": `^branch2$`,
					"branch3": `^branch3$`,
				},
				Validation: map[string]BranchValidationConfiguration{
					"branch1": map[string]string{
						"(A)": "contains B",
						"(B)": "contains A",
					},
					"branch2": map[string]string{
						"([AB])": "contains A or B (one or multiple)",
					},
					"branch3": map[string]string{},
				},
			}

			err := cfg.Validate(testData.BranchName, testData.CommitMessage)
			if testData.ErrorContains != "" {
				assert.Contains(t, err.Error(), testData.ErrorContains)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func createTestConfiguration() Configuration {
	cfg := Configuration{
		"test1": Project{Path: "/home"},
		"test2": Project{Path: "/home/project123"},
		"test3": Project{Path: ""},
	}
	return cfg
}
