package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	config, err := parse("test-data.yaml")
	if err != nil {
		t.Fatalf("Did not expect parse to return an error, but got: %v ", err)
	}

	expectedConfig := &Configuration{
		"project xyz": {
			Path: "/home/nils/projects/xyz/.git",
			BranchTypes: map[string]BranchTypePattern{
				"master":  `^(origin\/)*master`,
				"feature": `^(origin\/)*feature/.*$`,
				"develop": `^(origin\/)*develop$`,
				"release": `^(origin\/)*release\/v([0-9]*\.*)*$`,
				"hotfix":  `^(origin\/)*hotfix\/v([0-9]*\.*)*$`,
			},
			Templates: map[string]BranchTypeTemplate{
				"*": `{.BranchName}: {.CommitMessage}`,
			},
			Validation: map[string]BranchValidationConfiguration{
				"*": {
					`(?m)(?:\s|^|/)(([A-Z](_)*)+-[0-9]+)([\s,;:!.-]|$)`: "valid ticket ID (fallback validator)",
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
	_, err := parse("test-data-invalid.yaml")
	assert.Contains(t, err.Error(), "yaml: unmarshal errors")
}

func TestParse_FileNotFound_ExpectFileNotFoundError(t *testing.T) {
	_, err := parse("test-data-which-does-not-exist.yaml")
	assert.Contains(t, err.Error(), "no such file or directory")
}
