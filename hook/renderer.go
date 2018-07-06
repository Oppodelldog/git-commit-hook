package hook

import (
	"bytes"
	"github.com/Oppodelldog/git-commit-hook/config"
	"text/template"
)

type CommitMessageRenderer struct {
	projConf config.Project
}

// RenderCommitMessage renders a commit message using the template defined for the given resolveBranchNameFunc
func (r *CommitMessageRenderer) RenderCommitMessage(branchName string, viewModel ViewModel) (string, error) {
	branchType := r.projConf.GetBranchType(branchName)
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

func (r *CommitMessageRenderer) getTemplate(branchType string) string {
	foundTemplate := ""
	for configBranchType, branchTypeTemplate := range r.projConf.Templates {
		if configBranchType == branchType || configBranchType == "*" && foundTemplate == "" {
			foundTemplate = string(branchTypeTemplate)
		}
	}

	return foundTemplate
}