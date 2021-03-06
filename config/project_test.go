package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBranchType(t *testing.T) {

	testDataSet := []struct{ BranchName, ExpectedBranchType string }{
		{BranchName: "feature/FEATURE-123", ExpectedBranchType: "feature"},
		{BranchName: "release/v1.0.0", ExpectedBranchType: "release"},
		{BranchName: "release/v1.0.0-fix", ExpectedBranchType: "release"},
		{BranchName: "master", ExpectedBranchType: ""},
		{BranchName: "develop", ExpectedBranchType: ""},
	}

	for k, testData := range testDataSet {
		t.Run(string(k), func(t *testing.T) {

			cfg := &Project{
				BranchTypes: map[string]BranchTypePattern{
					"feature": `^(origin\/)*feature/.*$`,
					"release": `^(origin\/)*release\/v([0-9]*\.*)*(-fix)*$`,
				},
			}

			branchType := cfg.GetBranchType(testData.BranchName)

			assert.Exactly(t, testData.ExpectedBranchType, branchType)
		})
	}
}
