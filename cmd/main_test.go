package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"reflect"

	"path"

	"github.com/Oppodelldog/git-commit-hook/config"
	"github.com/Oppodelldog/git-commit-hook/subcommand"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

const testPath = "/tmp/git-commit-hook"
const featureBranch = "feature/PROJECT-123"
const nonFeatureBranch = "release/v0.1.2"
const commitMessageFile = ".git/COMMIT_EDITMSG"

//originals contains original values might be overwritten in tests
var originals = struct {
	osArgs                   []string
	rewriteCommitMessageFunc rewriteCommitMessageFuncDef
	diagnosticsFunc          callWithIntResult
	testFunc                 callWithIntResult
	installFunc              callWithIntResult
	uninstallFunc            callWithIntResult
	exitFunc                 exitFuncDef
	osStdout                 *os.File
}{
	osArgs: os.Args,
	rewriteCommitMessageFunc: rewriteCommitMessageFunc,
	diagnosticsFunc:          diagnosticsFunc,
	testFunc:                 testFunc,
	installFunc:              installFunc,
	uninstallFunc:            uninstallFunc,
	exitFunc:                 exitFunc,
	osStdout:                 os.Stdout,
}

//restoreOriginals restores original values
func restoreOriginals() {
	diagnosticsFunc = originals.diagnosticsFunc
	testFunc = originals.testFunc
	installFunc = originals.installFunc
	uninstallFunc = originals.uninstallFunc
	os.Args = originals.osArgs
	rewriteCommitMessageFunc = originals.rewriteCommitMessageFunc
	exitFunc = originals.exitFunc
	os.Stdout = originals.osStdout
	os.RemoveAll(testPath)
}

func TestMain_HappyPath(t *testing.T) {
	defer restoreOriginals()

	initGitRepositoryWithBranchAndConfig(t, featureBranch)
	setCommitMessage(t, "initial commit")
	prepareGitHookCall()

	assertProgramExistsWith(t, 0)

	main()

	expectedCommitMessage := fmt.Sprintf("%s: initial commit", featureBranch)
	assertCommitMessage(t, expectedCommitMessage)
}

func TestMain_ConfigurationNotFound(t *testing.T) {
	defer restoreOriginals()

	initTestFolder(t)
	initGitRepository(t, featureBranch)

	initialCommitMessage := "we expect this not to be changed by the tool"
	setCommitMessage(t, initialCommitMessage)
	prepareGitHookCall()

	assertProgramExistsWith(t, 1)

	w, stdOutChannel := captureStdOut(t)

	main()

	w.Close()

	stdOutput := <-stdOutChannel
	assert.Contains(t, stdOutput, "could not find config file")
	assertCommitMessage(t, initialCommitMessage)
}

func readCommitMessage(t *testing.T) string {
	b, err := ioutil.ReadFile(commitMessageFile)
	if err != nil {
		t.Fatalf("Did not exepect ioutil.ReadFile to return an error, but got: %v", err)
	}

	return string(b)
}

func prepareGitHookCall() {
	os.Args = []string{"git", commitMessageFile}
}

func setCommitMessage(t *testing.T, commitMessage string) {
	err := ioutil.WriteFile(commitMessageFile, []byte(commitMessage), 0777)
	if err != nil {
		t.Fatalf("Did not exepect ioutil.WriteFile to return an error, but got: %v", err)
	}
}

func TestMain_ErrorCase_TooFewArguments(t *testing.T) {
	defer restoreOriginals()

	testDataSet := map[string]struct{ OsArgs []string }{
		"no cli argument":  {[]string{}},
		"one cli argument": {[]string{"fyi: in runtime this will be the filepath to the command itself"}},
	}

	for testCaseName, testData := range testDataSet {
		t.Run(testCaseName, func(t *testing.T) {
			os.Args = testData.OsArgs

			assertProgramExistsWith(t, 0)
			w, stdOutChannel := captureStdOut(t)

			main()

			w.Close()

			stdOutput := <-stdOutChannel

			assert.Contains(t, stdOutput, "too few arguments")
		})
	}
}

func TestMain_FirstArgIsSubCommand_AppropriateFuncCalled(t *testing.T) {
	defer restoreOriginals()

	testDataSet := map[string]struct {
		SubCommandName string
		expectedFunc   *callWithIntResult
	}{
		"test":      {"test", &testFunc},
		"install":   {"test", &installFunc},
		"uninstall": {"test", &uninstallFunc},
		"diag":      {"test", &diagnosticsFunc},
	}

	for subCommandName, testData := range testDataSet {
		t.Run(subCommandName, func(t *testing.T) {
			restoreOriginals()
			os.Args = []string{"", subCommandName}

			funcStubCalled := false
			*testData.expectedFunc = func() int {
				funcStubCalled = true
				return 0
			}
			exitFunc = func(int) {}

			main()

			assert.True(t, funcStubCalled)
		})
	}
}

func TestMain_ErrorCase_EmptyCommitMessageFileName(t *testing.T) {
	defer restoreOriginals()

	os.Args = []string{"not important", ""}
	w, stdOutChannel := captureStdOut(t)

	assertProgramExistsWith(t, 1)

	main()

	w.Close()

	stdOutput := <-stdOutChannel
	assert.Contains(t, stdOutput, "no commit message file passed as parameter 1")
}

func TestMain_ErrorCase_CommitMessageFileNotFound(t *testing.T) {
	defer restoreOriginals()
	initTestFolder(t)
	writeConfigFile(t, testPath)

	os.Args = []string{"git", commitMessageFile}
	w, stdOutChannel := captureStdOut(t)

	assertProgramExistsWith(t, 1)

	main()

	w.Close()

	stdOutput := <-stdOutChannel
	assert.Contains(t, stdOutput, "error reading commit message from '.git/COMMIT_EDITMSG':")
}

func TestMain_ErrorCase_GitError(t *testing.T) {
	defer restoreOriginals()

	initGitRepositoryWithBranchAndConfig(t, nonFeatureBranch)
	setCommitMessage(t, "@noissue rc-fix")
	os.Args = []string{"git", commitMessageFile}
	w, stdOutChannel := captureStdOut(t)

	assertProgramExistsWith(t, 1)

	main()

	w.Close()

	stdOutput := <-stdOutChannel
	assert.Exactly(t, "error modifying commit message: validation error for branch 'release/v0.1.2'", stdOutput)
}

func TestMain_ExitFuncUsesAppropriateOsFunc(t *testing.T) {
	assert.Exactly(t, reflect.ValueOf(os.Exit).Pointer(), reflect.ValueOf(exitFunc).Pointer())
}

func TestMain_DiagnoseFuncMappedCorrectly(t *testing.T) {
	assert.Exactly(t, reflect.ValueOf(subcommand.Diagnostics).Pointer(), reflect.ValueOf(diagnosticsFunc).Pointer())
}

func TestMain_TestFuncMappedCorrectly(t *testing.T) {
	assert.Exactly(t, reflect.ValueOf(subcommand.Test).Pointer(), reflect.ValueOf(testFunc).Pointer())
}

func TestMain_InstallFuncMappedCorrectly(t *testing.T) {
	assert.Exactly(t, reflect.ValueOf(subcommand.Install).Pointer(), reflect.ValueOf(installFunc).Pointer())
}

func TestMain_UnnstallFuncMappedCorrectly(t *testing.T) {
	assert.Exactly(t, reflect.ValueOf(subcommand.Uninstall).Pointer(), reflect.ValueOf(uninstallFunc).Pointer())
}
func assertProgramExistsWith(t *testing.T, expectedExitCode int) {
	exitFunc = func(exitCode int) {
		assert.Exactly(t, expectedExitCode, exitCode)
	}
}

func captureStdOut(t *testing.T) (*os.File, chan string) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Did not exepect os.Pipe to return an error, but got: %v", err)
	}
	os.Stdout = w
	stdOutChannel := make(chan string)
	go func() {
		defer r.Close()

		scanner := bufio.NewScanner(r)
		if scanner.Scan() {
			stdOutChannel <- scanner.Text()
		}

		close(stdOutChannel)
	}()
	return w, stdOutChannel
}

func initGitRepositoryWithBranchAndConfig(t *testing.T, branchName string) {
	initTestFolder(t)
	initGitRepository(t, branchName)
	writeConfigFile(t, path.Join(testPath, ".git"))
}

func initTestFolder(t *testing.T) {
	err := os.RemoveAll(testPath)
	if err != nil {
		t.Fatalf("Did not expect os.RemoveAll to return an error, but got: %v ", err)
	}

	err = os.MkdirAll(testPath, 0777)
	if err != nil {
		t.Fatalf("Did not expect os.MkdirAll to return an error, but got: %v ", err)
	}

	err = os.Chdir(testPath)
	if err != nil {
		t.Fatalf("Did not expect os.Chdir to return an error, but got: %v ", err)
	}
}

func initGitRepository(t *testing.T, branchName string) {
	git(t, "init")
	err := ioutil.WriteFile("README.md", []byte("# test file"), 0777)
	if err != nil {
		t.Fatalf("could not write README.md: %v", err)
	}
	git(t, "config", "user.email", "odog@git-commit-hook.ok")
	git(t, "config", "user.name", "odog")
	git(t, "add", "-A")
	git(t, "commit", "-m", "initial commit")
	git(t, "checkout", "-b", branchName)
}

func writeConfigFile(t *testing.T, dir string) {
	os.MkdirAll(dir, 0777)
	cfg := config.Configuration{
		"test project": config.Project{
			Path: "/tmp/git-commit-hook/.git",
			BranchTypes: map[string]config.BranchTypePattern{
				"feature": `^feature/PROJECT-123$`,
				"release": `^release.*$`,
			},
			Templates: map[string]config.BranchTypeTemplate{
				"feature": "{{.BranchName}}: {{.CommitMessage}}",
			},
			Validation: map[string]config.BranchValidationConfiguration{
				"release": {
					"(?m)(?:\\s|^|/)(([A-Z](_)*)+-[0-9]+)([\\s,;:!.-]|$)": "valid ticket ID",
				},
			},
		},
	}

	configBytes, err := yaml.Marshal(cfg)
	if err != nil {
		t.Fatalf("Did not expect yaml.Marshal to return an error, but got: %v ", err)
	}

	err = ioutil.WriteFile(path.Join(dir, "git-commit-hook.yaml"), configBytes, 0666)
	if err != nil {
		t.Fatalf("Did not expect ioutil.WriteFile to return an error, but got: %v ", err)
	}
}

func git(t *testing.T, args ...string) {
	o, err := exec.Command("git", args...).CombinedOutput()
	if err != nil {
		t.Fatalf("'git %v init' failed with error: %v - output: %s", args, err, string(o))
	}
}

func assertCommitMessage(t *testing.T, expectedCommitMessage string) {
	modifiedCommitMessage := readCommitMessage(t)
	assert.Exactly(t, expectedCommitMessage, modifiedCommitMessage)
}
