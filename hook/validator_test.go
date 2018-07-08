package hook

import (
	"testing"

	"github.com/Oppodelldog/git-commit-hook/config"
	"github.com/stretchr/testify/assert"
)

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
		{
			BranchName:    "exotic branch name",
			CommitMessage: "hello",
			ErrorContains: "any branch without type falls into this validation rule",
		},
	}

	for testName, testData := range testDataSet {
		t.Run(string(testName), func(t *testing.T) {

			cfg := config.Project{
				BranchTypes: map[string]config.BranchTypePattern{
					"branch1": `^branch1$`,
					"branch2": `^branch2$`,
					"branch3": `^branch3$`,
				},
				Validation: map[string]config.BranchValidationConfiguration{
					"branch1": map[string]string{
						"(A)": "contains B",
						"(B)": "contains A",
					},
					"branch2": map[string]string{
						"([AB])": "contains A or B (one or multiple)",
					},
					"branch3": map[string]string{},
					"*": map[string]string{
						"^FALLBACKRULE$": "any branch without type falls into this validation rule",
					},
				},
			}

			validator := NewCommitMessageValidator(cfg)
			err := validator.Validate(testData.BranchName, testData.CommitMessage)
			if testData.ErrorContains != "" {
				assert.Contains(t, err.Error(), testData.ErrorContains)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
