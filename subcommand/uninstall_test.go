package subcommand

import (
	"testing"

	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/Oppodelldog/git-commit-hook/config"
	"github.com/Oppodelldog/git-commit-hook/testhelper"
	"github.com/stretchr/testify/assert"
)

func TestUninstallCommand_Uninstall_HappyPath(t *testing.T) {
	originArgs := os.Args
	defer func() { os.Args = originArgs }()
	defer testhelper.CleanupTestEnvironment(t)
	testhelper.PreapreTestEnvironment(t)
	testhelper.InitGitRepository(t, "feature/123")

	os.Args = []string{"git-commit-hook", "uninstall", "-p", "projectA"}

	configWithAProject := &config.Configuration{"projectA": config.Project{Path: testhelper.TestPathGitFolder}}

	existingFilePath := path.Join(testhelper.TestPathHooksFolder, gitCommitMessageHookName)
	err := ioutil.WriteFile(existingFilePath, []byte("TEST"), 0666)
	if err != nil {
		t.Fatalf("Did not expect ioutil.WriteFile to return an error, but got: %v ", err)
	}

	if _, err := os.Stat(existingFilePath); os.IsNotExist(err) {
		t.Fatalf("hook test-file missing")
	}

	cmd := NewUninstallerCommand()
	cmd.loadConfiguration = func() (*config.Configuration, error) { return configWithAProject, nil }
	res := cmd.Uninstall()

	assert.Exactly(t, 0, res)
	if _, err := os.Stat(existingFilePath); os.IsExist(err) {
		t.Fatalf("hook test-file not removed")
	}

}
func TestUninstallCommand_Uninstall_CanLoadConfiguration_ButNoParameterGiven_ShowsUsage(t *testing.T) {
	originArgs := os.Args
	os.Args = []string{"programm name", "install"}
	defer func() { os.Args = originArgs }()

	cmd := NewUninstallerCommand()
	cmd.stdoutWriter = bytes.NewBufferString("")
	cmd.loadConfiguration = func() (*config.Configuration, error) {
		return nil, nil
	}

	res := cmd.Uninstall()

	expectedOutput := `
Usage of git-commit-hook uninstall:
  -a	all
  -p string
    	project name
`
	assert.Exactly(t, 0, res)
	assert.Exactly(t, strings.TrimLeft(expectedOutput, "\n"), cmd.stdoutWriter.(*bytes.Buffer).String())
}

func TestUninstallCommand_Uninstall_UnknownParameter_ShowError(t *testing.T) {
	originArgs := os.Args
	os.Args = []string{"programm name", "uninstall", "--unknown-parameter"}
	defer func() { os.Args = originArgs }()
	defer testhelper.CleanupTestEnvironment(t)
	testhelper.PreapreTestEnvironment(t)

	uninstaller := NewUninstallerCommand()
	uninstaller.stdoutWriter = bytes.NewBufferString("")
	res := uninstaller.Uninstall()

	expectedOutput := `
flag provided but not defined: -unknown-parameter
Usage of git-commit-hook uninstall:
  -a	all
  -p string
    	project name
`

	assert.Exactly(t, strings.TrimLeft(expectedOutput, "\n"), uninstaller.stdoutWriter.(*bytes.Buffer).String())
	assert.Exactly(t, 1, res)
}

func TestUninstallCommand_Uninstall_CannotLoadConfiguration_ShowsError(t *testing.T) {
	originArgs := os.Args
	os.Args = []string{"programm name", "install"}
	defer func() { os.Args = originArgs }()

	cmd := NewUninstallerCommand()
	cmd.stdoutWriter = bytes.NewBufferString("")
	cmd.loadConfiguration = func() (*config.Configuration, error) {
		return nil, errors.New("could not load configuration")
	}

	res := cmd.Uninstall()

	expectedOutput := `
could not load configuration
`
	assert.Exactly(t, 1, res)
	assert.Exactly(t, strings.TrimLeft(expectedOutput, "\n"), cmd.stdoutWriter.(*bytes.Buffer).String())
}

func TestUninstallCommand_Uninstall_UninstallForAllAndOneProject(t *testing.T) {
	originArgs := os.Args
	defer func() { os.Args = originArgs }()
	configWithTwoProjects := &config.Configuration{"projectA": config.Project{Path: "pathA"}, "projectB": config.Project{Path: "pathB"}}

	testDataSet := map[string]struct {
		configuration    *config.Configuration
		additionalOsArgs []string
		expectedOutput   string
	}{
		"all": {
			additionalOsArgs: []string{"-a"},
			configuration:    configWithTwoProjects,
			expectedOutput: `
uninstalling git-commit-hook from 'pathA/hooks/commit-msg': OK
uninstalling git-commit-hook from 'pathB/hooks/commit-msg': OK
`,
		},
		"for one project": {
			additionalOsArgs: []string{"-p", "projectB"},
			configuration:    configWithTwoProjects,
			expectedOutput:   "uninstalling git-commit-hook from 'pathB/hooks/commit-msg': OK\n",
		},
	}

	for testCaseName, testData := range testDataSet {
		t.Run(testCaseName, func(t *testing.T) {

			os.Args = []string{"programm name", "uninstall"}
			os.Args = append(os.Args, testData.additionalOsArgs...)

			cmd := NewUninstallerCommand()
			cmd.stdoutWriter = bytes.NewBufferString("")
			cmd.loadConfiguration = func() (*config.Configuration, error) {
				return testData.configuration, nil
			}
			cmd.deleteFile = func(string) error { return nil }

			res := cmd.Uninstall()

			assert.Exactly(t, 0, res)

			assertContainsEachLine(t, testData.expectedOutput, cmd.stdoutWriter.(*bytes.Buffer).String())
		})
	}
}

func TestUninstallCommand_Uninstall_UninstallForAllAndOneProject_DeleteFileReturnsError(t *testing.T) {
	originArgs := os.Args
	defer func() { os.Args = originArgs }()
	configWithTwoProjects := &config.Configuration{"projectA": config.Project{Path: "pathA"}, "projectB": config.Project{Path: "pathB"}}

	testDataSet := map[string]struct {
		configuration    *config.Configuration
		additionalOsArgs []string
		expectedOutput   string
	}{
		"all": {
			additionalOsArgs: []string{"-a"},
			configuration:    configWithTwoProjects,
			expectedOutput: `
uninstalling git-commit-hook from 'pathA/hooks/commit-msg': could not delete file
uninstalling git-commit-hook from 'pathB/hooks/commit-msg': could not delete file
done with errors
`,
		},
		"for one project": {
			additionalOsArgs: []string{"-p", "projectB"},
			configuration:    configWithTwoProjects,
			expectedOutput:   "uninstalling git-commit-hook from 'pathB/hooks/commit-msg': could not delete file\n",
		},
	}

	for testCaseName, testData := range testDataSet {
		t.Run(testCaseName, func(t *testing.T) {

			os.Args = []string{"programm name", "uninstall"}
			os.Args = append(os.Args, testData.additionalOsArgs...)

			cmd := NewUninstallerCommand()
			cmd.stdoutWriter = bytes.NewBufferString("")
			cmd.loadConfiguration = func() (*config.Configuration, error) {
				return testData.configuration, nil
			}

			cmd.deleteFile = func(string) error { return errors.New("could not delete file") }

			res := cmd.Uninstall()

			assert.Exactly(t, 1, res)
			assertContainsEachLine(t, testData.expectedOutput, cmd.stdoutWriter.(*bytes.Buffer).String())
		})
	}
}

func TestUninstallCommand_Uninstall_ProjectNameCannotBeFound_ShowError(t *testing.T) {

	originArgs := os.Args
	os.Args = []string{"programm name", "install", "-p", "projectA"}
	defer func() { os.Args = originArgs }()

	cmd := NewUninstallerCommand()
	cmd.stdoutWriter = bytes.NewBufferString("")
	cmd.loadConfiguration = func() (*config.Configuration, error) {
		return &config.Configuration{}, nil
	}
	res := cmd.Uninstall()

	expectedOutput := `
project configuration not found for project name 'projectA'
`
	assert.Exactly(t, 1, res)
	assert.Exactly(t, strings.TrimLeft(expectedOutput, "\n"), cmd.stdoutWriter.(*bytes.Buffer).String())

}
