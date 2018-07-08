package hook

import (
	"bytes"
	"text/template"

	"github.com/Oppodelldog/git-commit-hook/config"
)

// CommitMessageRenderer implements rendering a commit message by a given
type CommitMessageRenderer interface {
	Render(viewModel ViewModel) (string, error)
}

//NewCommitMessageRenderer create a new CommitMessageRenderer
func NewCommitMessageRenderer(projConf config.Project) CommitMessageRenderer {
	return &commitMessageRenderer{
		projConf: projConf,
	}
}

// commitMessageRenderer renders a commit message by resolving the branch names template from project configuration
type commitMessageRenderer struct {
	projConf config.Project
}

// Render renders a commit message using the template defined for the given resolveBranchNameFunc
func (r *commitMessageRenderer) Render(viewModel ViewModel) (string, error) {
	branchType := r.projConf.GetBranchType(viewModel.BranchName)
	commitMessageTemplate := r.getTemplate(branchType)
	if commitMessageTemplate == "" {
		commitMessageTemplate = getFallbackCommitMessageTemplate()
	}
	tmpl, err := template.New("commitMessageTemplate").Parse(commitMessageTemplate)
	if err != nil {
		return "", err
	}
	buffer := bytes.NewBufferString("")
	err = tmpl.Execute(buffer, viewModel)

	return buffer.String(), err
}

func getFallbackCommitMessageTemplate() string {
	return "{{.CommitMessage}}"
}

func (r *commitMessageRenderer) getTemplate(branchType string) string {
	foundTemplate := ""
	for configBranchType, branchTypeTemplate := range r.projConf.Templates {
		if configBranchType == branchType || configBranchType == "*" && foundTemplate == "" {
			foundTemplate = string(branchTypeTemplate)
		}
	}

	return foundTemplate
}
