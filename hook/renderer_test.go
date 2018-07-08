package hook

import (
	"testing"

	"github.com/Oppodelldog/git-commit-hook/config"
	"github.com/stretchr/testify/assert"
)

func TestRenderCommitMessage(t *testing.T) {
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
	renderer := commitMessageRenderer{*cfg}
	modifiedCommitMessage, err := renderer.Render(viewModel)
	if err != nil {
		t.Fatalf("Did not expect Render to return an error, but got: %v ", err)
	}

	expectedCommitMessage := "feature/PRJ_TEST-1242: initial commit"
	assert.Exactly(t, expectedCommitMessage, modifiedCommitMessage)
}

func TestRenderCommitMessage_InvalidTemplate_ReturnsError(t *testing.T) {
	invalidTemplate := "{{{{{ HELLO"
	viewModel := ViewModel{}
	cfg := &config.Project{
		BranchTypes: map[string]config.BranchTypePattern{
			"feature": `(?m)^()((?!master|release|develop).)*$`,
		},
		Templates: map[string]config.BranchTypeTemplate{
			"feature": config.BranchTypeTemplate(invalidTemplate),
		}}

	renderer := commitMessageRenderer{*cfg}
	_, err := renderer.Render(viewModel)

	assert.Contains(t, err.Error(), "template:")
}

func TestRenderCommitMessage_NoTemplateFound_PassesBackTheGivenCommitMessage(t *testing.T) {
	givenCommitMessage := "some commit message"
	viewModel := ViewModel{
		CommitMessage: givenCommitMessage,
	}
	cfg := &config.Project{
		BranchTypes: map[string]config.BranchTypePattern{
			"feature": `(?m)^()((?!master|release|develop).)*$`,
		},
		Templates: map[string]config.BranchTypeTemplate{}}

	renderer := commitMessageRenderer{*cfg}
	modifiedCommitMessage, err := renderer.Render(viewModel)
	if err != nil {
		t.Fatalf("Did not expect Render to return an error, but got: %v ", err)
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
	renderer := commitMessageRenderer{*cfg}
	template := renderer.getTemplate("branch2")
	assert.Exactly(t, "templ2", template)

	template = renderer.getTemplate("branch4")
	assert.Exactly(t, "templ4", template)

	template = renderer.getTemplate("branch0")
	assert.Exactly(t, "fallback", template)
}

func TestRenderCommitMessage_NoBranchConfiguration_PassesBackTheGivenCommitMessage(t *testing.T) {
	givenCommitMessage := "some commit message"
	viewModel := ViewModel{
		CommitMessage: givenCommitMessage,
	}
	cfg := &config.Project{
		BranchTypes: map[string]config.BranchTypePattern{},
		Templates:   map[string]config.BranchTypeTemplate{}}

	renderer := commitMessageRenderer{*cfg}
	modifiedCommitMessage, err := renderer.Render(viewModel)
	if err != nil {
		t.Fatalf("Did not expect Render to return an error, but got: %v ", err)
	}

	assert.Exactly(t, givenCommitMessage, modifiedCommitMessage)
}
