package subcommand

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/Oppodelldog/git-commit-hook/testhelper"
	"github.com/stretchr/testify/assert"
	"os"
	"github.com/Oppodelldog/git-commit-hook/config"
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

func TestIsCommitHookInstalled(t *testing.T) {
	defer testhelper.CleanupTestEnvironment(t)

	testDataSet := map[string]struct {
		expectedResult bool
		prepareTest    func()
	}{
		"hook already installed": {
			prepareTest: func() {
				execFilePath, err := os.Executable()
				if err != nil {
					t.Fatalf("Did not expect os.Executable to return an error, but got: %v ", err)
				}
				os.Symlink(execFilePath, path.Join(testhelper.TestPathHooksFolder, gitCommitMessageHookName))
			},
			expectedResult: true,
		},
		"hook installed, but pointing to some folder": {
			prepareTest: func() {
				os.Symlink(testhelper.TestPath, path.Join(testhelper.TestPathHooksFolder, gitCommitMessageHookName))
			},
			expectedResult: false,
		},
		"hook installed, but no symlink": {
			prepareTest: func() {
				ioutil.WriteFile(path.Join(testhelper.TestPathHooksFolder, gitCommitMessageHookName), []byte("HOOK"), 0666)
			},
			expectedResult: false,
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
			res := isCommitHookInstalled(testhelper.TestPathGitFolder)

			assert.Exactly(t, testData.expectedResult, res)
		})
	}
}

func TestLoadProjectConfiguration(t *testing.T) {
	defer testhelper.CleanupTestEnvironment(t)

	testDataSet := map[string]struct {
		expectedErrorMessage string
		prepare              func()
	}{
		"config exists": {
			expectedErrorMessage: "",
			prepare: func() {
				testhelper.InitTestFolder(t)
				testhelper.WriteConfigFile(t, testhelper.TestPath)
			},
		},
		"config not found": {
			expectedErrorMessage: "project configuration not found for path",
			prepare:              func() {
				testhelper.InitTestFolder(t)
			},
		},
		"config not found sind wd doees not exist": {
			expectedErrorMessage: "getwd: no such file or directory",
			prepare:              func() {},
		},
	}

	for testCaseName, testData := range testDataSet {
		t.Run(testCaseName, func(t *testing.T) {
			testhelper.CleanupTestEnvironment(t)

			testData.prepare()
			configuration, err := loadProjectConfiguration()

			assert.IsType(t, config.Project{}, configuration)
			if testData.expectedErrorMessage != "" {
				assert.Contains(t, err.Error(), testData.expectedErrorMessage)
			}else{
				assert.NoError(t, err)
			}

		})
	}
}
