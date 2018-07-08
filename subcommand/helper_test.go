package subcommand

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/Oppodelldog/git-commit-hook/testhelper"
	"github.com/stretchr/testify/assert"
)

func TestIsAnotherGitHookInstalled(t *testing.T) {
	defer testhelper.CleanupTestEnvironment(t)

	testDataSet := map[string]struct {
		expectedResult bool
		prepareTest    func()
	}{
		"hook already installed": {
			prepareTest: func() {
				ioutil.WriteFile(path.Join(testhelper.TestPathHooksFolder, gitCommitMessageHookName), []byte("HOOK"), 0666)
			},
			expectedResult: true,
		},
		"no hook installed": {
			prepareTest:    func() {},
			expectedResult: false,
		},
	}

	for testCaseName, testData := range testDataSet {
		t.Run(testCaseName, func(t *testing.T) {
			testhelper.InitTestFolder(t)
			testhelper.InitGitRepository(t, "develop")
			testData.prepareTest()
			res := isAnotherGitHookInstalled(testhelper.TestPathGitFolder)

			assert.Exactly(t, testData.expectedResult, res)
		})
	}
}
