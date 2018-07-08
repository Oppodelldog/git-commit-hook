package subcommand

import (
	"testing"

	"bytes"
	"errors"
	"os"
	"strings"

	"github.com/Oppodelldog/git-commit-hook/config"
	"github.com/stretchr/testify/assert"
)

func TestInstall_UnknownArgument_ShowsUsage(t *testing.T) {
	originArgs := os.Args
	os.Args = []string{"programm name", "install", "--unknown-argument"}
	defer func() { os.Args = originArgs }()

	cmd := NewInstallCommand()
	cmd.stdoutWriter = bytes.NewBufferString("")

	res := cmd.Install()

	expectedOutput := `
flag provided but not defined: -unknown-argument
Usage of git-commit-hook install:
  -a	all
  -f	force file creation by overwriting
  -p string
    	project name
`

	assert.Exactly(t, 1, res)
	assert.Exactly(t, strings.TrimLeft(expectedOutput, "\n"), cmd.stdoutWriter.(*bytes.Buffer).String())
}

func TestInstall_CannotLoadConfiguration_ShowsError(t *testing.T) {
	originArgs := os.Args
	os.Args = []string{"programm name", "install"}
	defer func() { os.Args = originArgs }()

	cmd := NewInstallCommand()
	cmd.stdoutWriter = bytes.NewBufferString("")
	cmd.loadConfiguration = func() (*config.Configuration, error) {
		return nil, errors.New("could not load configuration")
	}

	res := cmd.Install()

	expectedOutput := `
could not load configuration
`
	assert.Exactly(t, 1, res)
	assert.Exactly(t, strings.TrimLeft(expectedOutput, "\n"), cmd.stdoutWriter.(*bytes.Buffer).String())
}

func TestInstall_CanLoadConfiguration_ButNoParameterGiven_ShowsUsage(t *testing.T) {
	originArgs := os.Args
	os.Args = []string{"programm name", "install"}
	defer func() { os.Args = originArgs }()

	cmd := NewInstallCommand()
	cmd.stdoutWriter = bytes.NewBufferString("")
	cmd.loadConfiguration = func() (*config.Configuration, error) {
		return nil, nil
	}

	res := cmd.Install()

	expectedOutput := `
Usage of git-commit-hook install:
  -a	all
  -f	force file creation by overwriting
  -p string
    	project name
`
	assert.Exactly(t, 1, res)
	assert.Exactly(t, strings.TrimLeft(expectedOutput, "\n"), cmd.stdoutWriter.(*bytes.Buffer).String())
}

func TestInstall_InstallForAllProjects(t *testing.T) {
	originArgs := os.Args
	defer func() { os.Args = originArgs }()
	configWithTwoProjects := &config.Configuration{"projectA": config.Project{Path: "pathA"}, "projectB": config.Project{Path: "pathB"}}

	testDataSet := map[string]struct {
		gitHookInstallerExpectedParms []gitHookInstallerMockParams
		configuration                 *config.Configuration
		additionalOsArgs              []string
		expectedOutput                string
	}{
		"without -f": {
			configuration: configWithTwoProjects,
			gitHookInstallerExpectedParms: []gitHookInstallerMockParams{
				{"pathA", false},
				{"pathB", false},
			},
			expectedOutput: `
installing git-commit-hook to 'pathA': OK
installing git-commit-hook to 'pathB': OK
`,
		},
		"with -f": {
			additionalOsArgs: []string{"-f"},
			configuration:    configWithTwoProjects,
			gitHookInstallerExpectedParms: []gitHookInstallerMockParams{
				{"pathA", true},
				{"pathB", true},
			},
			expectedOutput: `
installing git-commit-hook to 'pathA': OK
installing git-commit-hook to 'pathB': OK
`,
		},
		"without -f and project name": {
			additionalOsArgs: []string{"-p", "projectB"},
			configuration:    configWithTwoProjects,
			gitHookInstallerExpectedParms: []gitHookInstallerMockParams{
				{"pathB", false},
			},
			expectedOutput: "installing git-commit-hook to 'pathB': OK\n",
		},
		"with -f and project name": {
			additionalOsArgs: []string{"-f", "-p", "projectA"},
			configuration:    configWithTwoProjects,
			gitHookInstallerExpectedParms: []gitHookInstallerMockParams{
				{"pathA", true},
			},
			expectedOutput: "installing git-commit-hook to 'pathA': OK\n",
		},
	}

	for testCaseName, testData := range testDataSet {
		t.Run(testCaseName, func(t *testing.T) {

			os.Args = []string{"programm name", "install", "-a"}
			os.Args = append(os.Args, testData.additionalOsArgs...)

			cmd := NewInstallCommand()
			cmd.stdoutWriter = bytes.NewBufferString("")
			cmd.loadConfiguration = func() (*config.Configuration, error) {
				return testData.configuration, nil
			}
			cmd.gitHookInstaller = &gitHookInstallerMock{
				t:      t,
				params: testData.gitHookInstallerExpectedParms,
			}

			res := cmd.Install()

			assert.Exactly(t, 0, res)

			assertContainsEachLine(t, testData.expectedOutput, cmd.stdoutWriter.(*bytes.Buffer).String())
		})
	}
}

func assertContainsEachLine(t *testing.T, expectedOutput string, output string) {
	expectedOutput = strings.TrimLeft(expectedOutput, "\n")
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		assert.Contains(t, output, line)
	}
}

func TestInstall_ProjectConfigCannotBeFound_ShowError(t *testing.T) {

	originArgs := os.Args
	os.Args = []string{"programm name", "install", "-p", "projectA"}
	defer func() { os.Args = originArgs }()

	cmd := NewInstallCommand()
	cmd.stdoutWriter = bytes.NewBufferString("")
	cmd.loadConfiguration = func() (*config.Configuration, error) {
		return &config.Configuration{}, nil
	}
	res := cmd.Install()

	expectedOutput := `
project configuration not found for project name 'projectA'
`
	assert.Exactly(t, 1, res)
	assert.Exactly(t, strings.TrimLeft(expectedOutput, "\n"), cmd.stdoutWriter.(*bytes.Buffer).String())

}

func TestInstall_GitHookInstallerReturnsError_ShowError(t *testing.T) {
	originArgs := os.Args
	defer func() { os.Args = originArgs }()

	const errorMessageStub = "some error in git hook installer"

	testDataSet := map[string]struct {
		additionalOsArgs []string
		expectedOutput   string
	}{
		"all":            {additionalOsArgs: []string{"-a"}, expectedOutput: "installing git-commit-hook to '': some error in git hook installer\ndone with errors\n"},
		"single project": {additionalOsArgs: []string{"-p", "projectA"}, expectedOutput: "installing git-commit-hook to '': " + errorMessageStub + "\n"},
	}
	for testCaseName, testData := range testDataSet {
		t.Run(testCaseName, func(t *testing.T) {

			os.Args = []string{"programm name", "install"}
			os.Args = append(os.Args, testData.additionalOsArgs...)

			cmd := NewInstallCommand()
			cmd.stdoutWriter = bytes.NewBufferString("")
			cmd.loadConfiguration = func() (*config.Configuration, error) {
				return &config.Configuration{"projectA": config.Project{}}, nil
			}
			cmd.gitHookInstaller = &gitHookInstallerErrorMock{errorMessageStub}
			res := cmd.Install()

			assert.Exactly(t, 1, res)
			assert.Exactly(t, testData.expectedOutput, cmd.stdoutWriter.(*bytes.Buffer).String())
		})
	}
}

type gitHookInstallerMockParams struct {
	p1 string
	p2 bool
}
type gitHookInstallerMock struct {
	t      *testing.T
	params []gitHookInstallerMockParams
	callNo int
}

func (m *gitHookInstallerMock) installForProject(gitFolderPath string, forceOverwrite bool) error {
	for k, v := range m.params {
		if v.p1 == gitFolderPath && v.p2 == forceOverwrite {
			m.params = append(m.params[:k], m.params[k+1:]...)
			return nil
		}
	}
	m.t.Fatalf("no call expectations given")
	return nil
}

type gitHookInstallerErrorMock struct {
	errorMessage string
}

func (m *gitHookInstallerErrorMock) installForProject(gitFolderPath string, forceOverwrite bool) error {
	return errors.New(m.errorMessage)
}
