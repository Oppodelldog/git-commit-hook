package subcommand

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitHookInstaller_Install(t *testing.T) {
	cleanupTestEnvironment(t)
	defer cleanupTestEnvironment(t)
	gitFolder := path.Join(testDir, ".git")
	gitHooksDir := path.Join(gitFolder, "hooks")
	gitHookFilePath := path.Join(gitHooksDir, gitCommitMessageHookName)
	err := os.MkdirAll(gitHooksDir, 0777)
	if err != nil {
		t.Fatalf("Did not expect os.MkdirAll to return an error, but got: %v ", err)
	}

	installer := NewGitHookInstaller()

	res := installer.installForProject(gitFolder, false)

	assertCurrentExecutableIsSymlinkedAsGitHook(t, gitHookFilePath)
	assert.NoError(t, res)
}

func TestGitHookInstaller_FileAlreadyExists_NoForce(t *testing.T) {
	cleanupTestEnvironment(t)
	defer cleanupTestEnvironment(t)
	gitFolder := path.Join(testDir, ".git")
	gitHooksDir := path.Join(gitFolder, "hooks")
	gitHookFilePath := path.Join(gitHooksDir, gitCommitMessageHookName)
	err := os.MkdirAll(gitHooksDir, 0777)
	if err != nil {
		t.Fatalf("Did not expect os.MkdirAll to return an error, but got: %v ", err)
	}

	ioutil.WriteFile(gitHookFilePath, []byte("TEST"), 0666)

	installer := NewGitHookInstaller()

	err = installer.installForProject(gitFolder, false)

	assert.Exactly(t, "file already exists, use -f to force overwriting", err.Error())
}

func TestGitHookInstaller_FileAlreadyExists_WithForce(t *testing.T) {
	cleanupTestEnvironment(t)
	defer cleanupTestEnvironment(t)
	gitFolder := path.Join(testDir, ".git")
	gitHooksDir := path.Join(gitFolder, "hooks")
	gitHookFilePath := path.Join(gitHooksDir, gitCommitMessageHookName)
	err := os.MkdirAll(gitHooksDir, 0777)
	if err != nil {
		t.Fatalf("Did not expect os.MkdirAll to return an error, but got: %v ", err)
	}

	ioutil.WriteFile(gitHookFilePath, []byte("TEST"), 0666)

	installer := NewGitHookInstaller()

	err = installer.installForProject(gitFolder, true)

	assert.NoError(t, err)
	assertCurrentExecutableIsSymlinkedAsGitHook(t, gitHookFilePath)
}

func TestGitHookInstaller_FileAlreadyExists_WithForce_RemoveFails(t *testing.T) {
	cleanupTestEnvironment(t)
	defer cleanupTestEnvironment(t)
	gitFolder := path.Join(testDir, ".git")
	gitHooksDir := path.Join(gitFolder, "hooks")
	gitHookFilePath := path.Join(gitHooksDir, gitCommitMessageHookName)
	err := os.MkdirAll(gitHooksDir, 0777)
	if err != nil {
		t.Fatalf("Did not expect os.MkdirAll to return an error, but got: %v ", err)
	}

	ioutil.WriteFile(gitHookFilePath, []byte("TEST"), 0666)

	installer := NewGitHookInstaller()

	removeFileErrorStub := errors.New("some error while removing file")
	installer.(*gitHookInstaller).removeFile = func(string) error {
		return removeFileErrorStub
	}

	err = installer.installForProject(gitFolder, true)

	assert.Exactly(t, removeFileErrorStub, err)
}

func TestGitHookInstaller_CreateSymlinkFails(t *testing.T) {
	cleanupTestEnvironment(t)
	defer cleanupTestEnvironment(t)
	gitFolder := path.Join(testDir, ".git")
	gitHooksDir := path.Join(gitFolder, "hooks")
	err := os.MkdirAll(gitHooksDir, 0777)
	if err != nil {
		t.Fatalf("Did not expect os.MkdirAll to return an error, but got: %v ", err)
	}
	installer := NewGitHookInstaller()

	createSymlinkErrorStub := errors.New("some error while removing file")
	installer.(*gitHookInstaller).createSymlink = func(p1 string, p2 string) error {
		return createSymlinkErrorStub
	}

	err = installer.installForProject(gitFolder, true)

	assert.Exactly(t, createSymlinkErrorStub, err)
}

func TestGitHookInstaller_CannotGetCurrentExecutableFilePath(t *testing.T) {
	cleanupTestEnvironment(t)
	defer cleanupTestEnvironment(t)
	gitFolder := path.Join(testDir, ".git")
	gitHooksDir := path.Join(gitFolder, "hooks")
	err := os.MkdirAll(gitHooksDir, 0777)
	if err != nil {
		t.Fatalf("Did not expect os.MkdirAll to return an error, but got: %v ", err)
	}
	installer := NewGitHookInstaller()

	getExecutableFilePathErrorStub := errors.New("some error while removing file")
	installer.(*gitHookInstaller).getCurrentExecutableFilePath = func() (string, error) {
		return "", getExecutableFilePathErrorStub
	}

	err = installer.installForProject(gitFolder, true)

	assert.Exactly(t, getExecutableFilePathErrorStub, err)
}

func assertCurrentExecutableIsSymlinkedAsGitHook(t *testing.T, gitHookFilePath string) {
	exeFile, err := os.Executable()
	if err != nil {
		t.Fatalf("Did not expect os.Executable to return an error, but got: %v ", err)
	}
	expectedFileBytes, err := ioutil.ReadFile(exeFile)
	if err != nil {
		t.Fatalf("Did not expect ioutil.ReadFile to return an error, but got: %v ", err)
	}
	fileBytes, err := ioutil.ReadFile(gitHookFilePath)
	if err != nil {
		t.Fatalf("Did not expect ioutil.ReadFile to return an error, but got: %v ", err)
	}
	assert.Exactly(t, expectedFileBytes, fileBytes)
}
