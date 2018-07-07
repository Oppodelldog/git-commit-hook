package hook

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/Oppodelldog/git-commit-hook/config"
)

func TestRenderCommitMessage(t *testing.T) {
	configBranchName := "feature"
	viewModel := ViewModel{
		BranchName:    "feature/PRJ_TEST-1242",
		CommitMessage: "initial commit",
	}
	cfg := &config.Project{
		BranchTypes: map[string]config.BranchTypePattern{
			"feature": `(?m)^()((?!master|release|develop).)*$`,
		},
		Templates: map[string]config.BranchTypeTemplate{
			"feature": "{{.BranchName}}: {{.CommitMessage}}",
		}}
	renderer := CommitMessageRenderer{*cfg}
	modifiedCommitMessage, err := renderer.RenderCommitMessage(configBranchName, viewModel)
	if err != nil {
		t.Fatalf("Did not expect RenderCommitMessage to return an error, but got: %v ", err)
	}

	expectedCommitMessage := "feature/PRJ_TEST-1242: initial commit"
	assert.Exactly(t, expectedCommitMessage, modifiedCommitMessage)
}

func TestRenderCommitMessage_InvalidTemplate_ReturnsError(t *testing.T) {
	configBranchName := "feature"
	invalidTemplate := "{{{{{ HELLO"
	viewModel := ViewModel{}
	cfg := &config.Project{
		BranchTypes: map[string]config.BranchTypePattern{
			"feature": `(?m)^()((?!master|release|develop).)*$`,
		},
		Templates: map[string]config.BranchTypeTemplate{
			"feature": config.BranchTypeTemplate(invalidTemplate),
		}}

	renderer := CommitMessageRenderer{*cfg}
	_, err := renderer.RenderCommitMessage(configBranchName, viewModel)

	assert.Contains(t, err.Error(), "template:")
}

func TestRenderCommitMessage_NoTemplateFound_PassesBackTheGivenCommitMessage(t *testing.T) {
	givenCommitMessage := "some commit message"
	configBranchName := "feature"
	viewModel := ViewModel{
		CommitMessage: givenCommitMessage,
	}
	cfg := &config.Project{
		BranchTypes: map[string]config.BranchTypePattern{
			"feature": `(?m)^()((?!master|release|develop).)*$`,
		},
		Templates: map[string]config.BranchTypeTemplate{}}

	renderer := CommitMessageRenderer{*cfg}
	modifiedCommitMessage, err := renderer.RenderCommitMessage(configBranchName, viewModel)
	if err != nil {
		t.Fatalf("Did not expect RenderCommitMessage to return an error, but got: %v ", err)
	}

	assert.Exactly(t, givenCommitMessage, modifiedCommitMessage)
}

func TestGetTemplate(t *testing.T) {
	cfg := &config.Project{
		Templates: map[string]config.BranchTypeTemplate{
			"branch1": "templ1",
			"branch2": "templ2",
			"*":       "fallback",
			"branch4": "templ4",
		},
	}
	renderer := CommitMessageRenderer{*cfg}
	template := renderer.getTemplate("branch2")
	assert.Exactly(t, "templ2", template)

	template = renderer.getTemplate("branch4")
	assert.Exactly(t, "templ4", template)

	template = renderer.getTemplate("branch0")
	assert.Exactly(t, "fallback", template)
}

func TestRenderCommitMessage_NoBranchConfiguration_PassesBackTheGivenCommitMessage(t *testing.T) {
	givenCommitMessage := "some commit message"
	configBranchName := "feature"
	viewModel := ViewModel{
		CommitMessage: givenCommitMessage,
	}
	cfg := &config.Project{
		BranchTypes: map[string]config.BranchTypePattern{},
		Templates:   map[string]config.BranchTypeTemplate{}}

	renderer := CommitMessageRenderer{*cfg}
	modifiedCommitMessage, err := renderer.RenderCommitMessage(configBranchName, viewModel)
	if err != nil {
		t.Fatalf("Did not expect RenderCommitMessage to return an error, but got: %v ", err)
	}

	assert.Exactly(t, givenCommitMessage, modifiedCommitMessage)
}
