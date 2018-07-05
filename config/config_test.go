package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBranchType(t *testing.T) {

	testDataSet := []struct{ BranchName, ExpectedBranchType string }{
		{BranchName: "FEATURE-123", ExpectedBranchType: "feature"},
		{BranchName: "noissue_fix_something", ExpectedBranchType: "feature"},
		{BranchName: "release/v1.0.0", ExpectedBranchType: "release"},
		{BranchName: "release/v1.0.0-fix", ExpectedBranchType: "release"},
		{BranchName: "master", ExpectedBranchType: ""},
		{BranchName: "develop", ExpectedBranchType: ""},
	}

	for k, testData := range testDataSet {
		t.Run(string(k), func(t *testing.T) {

			cfg := &ProjectConfiguration{
				BranchTypes: map[string]BranchTypePattern{
					"feature": `(?m)^()((?!master|release|develop).)*$`,
					"release": `(?m)^(origin\/)*release\/v([0-9]*\.*)*(-fix)*$`,
				},
			}

			branchType := cfg.GetBranchType(testData.BranchName)

			assert.Exactly(t, testData.ExpectedBranchType, branchType)
		})
	}
}

func TestGetProjectConfiguration(t *testing.T) {
	expectedProjectCfg := ProjectConfiguration{Path: "/home/project123"}
	cfg := Configuration{
		"test1": ProjectConfiguration{Path: "/home"},
		"test2": expectedProjectCfg,
		"test3": ProjectConfiguration{Path: ""},
	}
	path := "/home/project123"
	projectCfg, err := (&cfg).GetProjectConfiguration(path)
	if err != nil {
		t.Fatalf("Did not expect GetProjectConfiguration to return an error, but got: %v ", err)
	}

	assert.Exactly(t, expectedProjectCfg, projectCfg)
}

func TestGetProjectConfiguration_NotFound(t *testing.T) {
	cfg := Configuration{
		"test1": ProjectConfiguration{Path: "/home"},
		"test2": ProjectConfiguration{Path: ""},
	}
	path := "/home/does-not-exist"
	_, err := (&cfg).GetProjectConfiguration(path)

	assert.Contains(t, err.Error(), "not found")

}

func TestRenderCommitMessage(t *testing.T) {
	configBranchName := "feature"
	viewModel := ViewModel{
		BranchName:    "feature/PRJ_TEST-1242",
		CommitMessage: "initial commit",
	}
	cfg := &ProjectConfiguration{
		BranchTypes: map[string]BranchTypePattern{
			"feature": `(?m)^()((?!master|release|develop).)*$`,
		},
		Templates: map[string]BranchTypeTemplate{
			"feature": "{{.BranchName}}: {{.CommitMessage}}",
		}}

	modifiedCommitMessage, err := cfg.RenderCommitMessage(configBranchName, viewModel)
	if err != nil {
		t.Fatalf("Did not expect RenderCommitMessage to return an error, but got: %v ", err)
	}

	expectedCommitMessage := "feature/PRJ_TEST-1242: initial commit"
	assert.Exactly(t, expectedCommitMessage, modifiedCommitMessage)
}

func TestRenderCommitMessage_InvalidTemplate_ReturnsError(t *testing.T) {
	configBranchName := "feature"
	viewModel := ViewModel{}
	cfg := &ProjectConfiguration{
		BranchTypes: map[string]BranchTypePattern{
			"feature": `(?m)^()((?!master|release|develop).)*$`,
		},
		Templates: map[string]BranchTypeTemplate{
			"feature": "{{{{{ HELLO",
		}}

	_, err := cfg.RenderCommitMessage(configBranchName, viewModel)

	assert.Contains(t, err.Error(), "template:")
}

func TestGetTemplate(t *testing.T) {
	cfg := &ProjectConfiguration{
		Templates: map[string]BranchTypeTemplate{
			"branch1": "templ1",
			"branch2": "templ2",
			"*":       "fallback",
			"branch4": "templ4",
		},
	}

	template := cfg.getTemplate("branch2")
	assert.Exactly(t, "templ2", template)

	template = cfg.getTemplate("branch4")
	assert.Exactly(t, "templ4", template)

	template = cfg.getTemplate("branch0")
	assert.Exactly(t, "fallback", template)
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

			cfg := &ProjectConfiguration{
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
