package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	config, err := Parse("test-data.yaml")
	if err != nil {
		t.Fatalf("Did not expect Parse to return an error, but got: %v ", err)
	}

	expectedConfig := &Configuration{
		"project xyz": {
			Path: "/home/nils/projects/xyz",
			BranchTypes: map[string]BranchTypeConfiguration{
				"master":  {Pattern: `^(origin\/)*master`},
				"feature": {Pattern: `^((?!master|release|develop).)*$`},
				"develop": {Pattern: `^(origin\/)*develop$`},
				"release": {Pattern: `^(origin\/)*release\/v([0-9]*\.*)*$`},
				"hotfix":  {Pattern: `^(origin\/)*hotfix\/v([0-9]*\.*)*$`},
			},
			Templates: map[string]BranchTemplateConfiguration{
				"*": {Template: `{.BranchName}: {.CommitMessage}`},
			},
			Validation: map[string]BranchValidationConfiguration{
				"feature": {
					`(?m)(?:\s|^|/)(([A-Z](_)*)+-[0-9]+)([\s,;:!.-]|$)`: "valid ticket ID",
				},
				"develop": {
					`(?m)(?:\s|^|/)(([A-Z](_)*)+-[0-9]+)([\s,;:!.-]|$)`: "valid ticket ID",
					"(?m)@noissue":                                      "@noissue",
				},
				"release": {
					`(?m)(?:\s|^|/)(([A-Z](_)*)+-[0-9]+)([\s,;:!.-]|$)`: "valid ticket ID",
					"(?m)@rc-fix":                                       "an @rc-fix indicator",
					"(?m)@noissue":                                      "@noissue",
				},
				"master": {
					`(?m)(?:\s|^|/)(([A-Z](_)*)+-[0-9]+)([\s,;:!.-]|$)`: "valid ticket ID",
					"(?m)@noissue":                                      "@noissue",
				},
			},
		},
	}

	assert.Exactly(t, expectedConfig, config)
}

func TestParse_InvalidTestData_ExpectUnmarshalError(t *testing.T) {
	_, err := Parse("test-data-invalid.yaml")
	assert.Contains(t, err.Error(), "yaml: unmarshal errors")
}

func TestParse_FileNotFound_ExpectFileNotFoundError(t *testing.T) {
	_, err := Parse("test-data-which-does-not-exist.yaml")
	assert.Contains(t, err.Error(), "no such file or directory")
}
