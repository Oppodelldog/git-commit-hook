package gitcommithook

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsFeatureBranch_Positive(t *testing.T) {

	testData := []struct {
		branchName     string
		expectedResult bool
	}{
		{"TOOL-1242", true},
		{"feature/TOOL-1242", true},
	}

	for _, test := range testData {
		t.Run(test.branchName, func(t *testing.T) {
			res := IsFeatureBranch(test.branchName)
			assert.Exactly(t, test.expectedResult, res)
		})
	}
}

func TestIsFeatureBranch_Negative(t *testing.T) {
	testData := []struct {
		branchName     string
		expectedResult bool
	}{
		{"develop", false},
		{"master", false},
		{"release", false},
		{"release/v0.1.0", false},
		{"release/v0.1.0-fix", false},
		{"hotfix/v0.1.0", false},
	}

	for _, test := range testData {
		t.Run(test.branchName, func(t *testing.T) {
			res := IsFeatureBranch(test.branchName)
			assert.Exactly(t, test.expectedResult, res)
		})
	}
}
