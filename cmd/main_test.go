package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"reflect"

	"github.com/stretchr/testify/assert"
)

const testPath = "/tmp/git-commit-hook"
const featureBranch = "feature/PROJECT-123"
const nonFeatureBranch = "release/v0.1.2"
const commitMessageFile = ".git/COMMIT_EDITMSG"

//originals contains original values might be overwritten in tests
var originals = struct {
	osArgs                   []string
	rewriteCommitMessageFunc rewriteCommitMessageFuncDef
	exitFunc                 exitFuncDef
	osStdout                 *os.File
}{
	osArgs: os.Args,
	rewriteCommitMessageFunc: rewriteCommitMessageFunc,
	exitFunc:                 exitFunc,
	osStdout:                 os.Stdout,
}

//restoreOriginals restores original values
func restoreOriginals() {
	os.Args = originals.osArgs
	rewriteCommitMessageFunc = originals.rewriteCommitMessageFunc
	exitFunc = originals.exitFunc
	os.Stdout = originals.osStdout
	os.RemoveAll(testPath)
}

func TestMain_HappyPath(t *testing.T) {
	defer restoreOriginals()

	initGitRepositoryWithBranch(t, featureBranch)
	setCommitMessage(t, "initial commit")
	prepareGitHookCall()

	programExistsWith(t, 0)

	main()

	expectedCommitMessage := fmt.Sprintf("%s: initial commit", featureBranch)
	modifiedCommitMessage := readCommitMessage(t)
	assert.Exactly(t, expectedCommitMessage, modifiedCommitMessage)
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

func TestMain_ErrorCase_CommitMessageFileNotFound(t *testing.T) {
	defer restoreOriginals()

	os.Args = []string{"git", commitMessageFile}
	w, stdOutChannel := captureStdOut(t)

	programExistsWith(t, 1)

	main()

	w.Close()

	stdOutput := <-stdOutChannel
	assert.Exactly(t, "error reading commit message from '.git/COMMIT_EDITMSG': open .git/COMMIT_EDITMSG: no such file or directory", stdOutput)
}

func TestMain_ErrorCase_GitError(t *testing.T) {
	defer restoreOriginals()

	initGitRepositoryWithBranch(t, nonFeatureBranch)
	setCommitMessage(t, "@noissue rc-fix")
	os.Args = []string{"git", commitMessageFile}
	w, stdOutChannel := captureStdOut(t)

	programExistsWith(t, 1)

	main()

	w.Close()

	stdOutput := <-stdOutChannel
	assert.Exactly(t, "error modifying commit message: feature reference is required in '@noissue rc-fix'", stdOutput)
}

func TestMain_ExitFuncUsesAppropriateOsFunc(t *testing.T) {
	assert.Exactly(t, reflect.ValueOf(os.Exit).Pointer(), reflect.ValueOf(exitFunc).Pointer())
}

func programExistsWith(t *testing.T, expectedExitCode int) {
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

func initGitRepositoryWithBranch(t *testing.T, branchName string) {
	os.RemoveAll(testPath)
	os.MkdirAll(testPath, 0777)
	os.Chdir(testPath)
	git(t, "init")
	err := ioutil.WriteFile("README.md", []byte("# test file"), 0777)
	if err != nil {
		t.Fatalf("could not write README.md: %v", err)
	}
	git(t, "add", "-A")
	git(t, "commit", "-m", "initial commit")
	git(t, "checkout", "-b", branchName)
}

func git(t *testing.T, args ...string) {
	o, err := exec.Command("git", args...).CombinedOutput()
	if err != nil {
		t.Fatalf("'git %v init' failed with error: %v - output: %s", args, err, string(o))
	}
}
