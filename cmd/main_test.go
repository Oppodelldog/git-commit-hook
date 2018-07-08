package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"reflect"

	"path"

	"github.com/Oppodelldog/git-commit-hook/subcommand"
	"github.com/Oppodelldog/git-commit-hook/testhelper"
	"github.com/stretchr/testify/assert"
)

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
	os.RemoveAll(testhelper.TestPath)
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
	defer testhelper.CleanupTestEnvironment(t)
	testhelper.CleanupTestEnvironment(t)
	testhelper.InitTestFolder(t)
	testhelper.InitGitRepository(t, featureBranch)

	initialCommitMessage := "we expect this not to be changed by the tool"
	setCommitMessage(t, initialCommitMessage)
	prepareGitHookCall()

	assertProgramExistsWith(t, 1)

	w, stdOutChannel := captureStdOut(t)

	main()

	w.Close()

	stdOutput := <-stdOutChannel
	assert.Contains(t, stdOutput, "could not find config file at")
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
	testhelper.InitTestFolder(t)
	testhelper.WriteConfigFile(t, testhelper.TestPath)

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
	assert.Exactly(t, reflect.ValueOf(subcommand.NewDiagCommand().Diagnostics).Pointer(), reflect.ValueOf(diagnosticsFunc).Pointer())
}

func TestMain_TestFuncMappedCorrectly(t *testing.T) {
	assert.Exactly(t, reflect.ValueOf(subcommand.NewTestCommand().Test).Pointer(), reflect.ValueOf(testFunc).Pointer())
}

func TestMain_InstallFuncMappedCorrectly(t *testing.T) {
	assert.Exactly(t, reflect.ValueOf(subcommand.NewInstallCommand().Install).Pointer(), reflect.ValueOf(installFunc).Pointer())
}

func TestMain_UnnstallFuncMappedCorrectly(t *testing.T) {
	assert.Exactly(t, reflect.ValueOf(subcommand.NewUninstallerCommand().Uninstall).Pointer(), reflect.ValueOf(uninstallFunc).Pointer())
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
	testhelper.InitTestFolder(t)
	testhelper.InitGitRepository(t, branchName)
	testhelper.WriteConfigFile(t, path.Join(testhelper.TestPath, ".git"))
}

func assertCommitMessage(t *testing.T, expectedCommitMessage string) {
	modifiedCommitMessage := readCommitMessage(t)
	assert.Exactly(t, expectedCommitMessage, modifiedCommitMessage)
}
