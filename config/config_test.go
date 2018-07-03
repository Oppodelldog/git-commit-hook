package config

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

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
		Templates: map[string]BranchTemplateConfiguration{
			"feature": {Template: "{{.BranchName}}: {{.CommitMessage}}"},
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
		Templates: map[string]BranchTemplateConfiguration{
			"feature": {Template: "{{{{{ HELLO"},
		}}

	_, err := cfg.RenderCommitMessage(configBranchName, viewModel)

	assert.Contains(t, err.Error(), "template:")
}

func TestGetTemplate(t *testing.T) {
	cfg := &ProjectConfiguration{
		Templates: map[string]BranchTemplateConfiguration{
			"branch1": {Template: "templ1"},
			"branch2": {Template: "templ2"},
			"*":       {Template: "fallback"},
			"branch4": {Template: "templ4"},
		},
	}

	template := cfg.GetTemplate("branch2")
	assert.Exactly(t, "templ2", template)

	template = cfg.GetTemplate("branch4")
	assert.Exactly(t, "templ4", template)

	template = cfg.GetTemplate("branch0")
	assert.Exactly(t, "fallback", template)
}
