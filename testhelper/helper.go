package testhelper

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/Oppodelldog/git-commit-hook/config"
	"gopkg.in/yaml.v2"
)

// TestPath holds the absolute path to the test folder
const TestPath = "/tmp/git-commit-hook"

// TestPathGitFolder holds the .git folder for the test
const TestPathGitFolder = "/tmp/git-commit-hook/.git"

// TestPathHooksFolder holds the hooks folder inside the .git folder for the test
const TestPathHooksFolder = "/tmp/git-commit-hook/.git/hooks"

//InitGitRepository initializes a test git repository with the given branch name in the current directory
func InitGitRepository(t *testing.T, branchName string) {
	t.Helper()
	Git(t, "init")
	err := ioutil.WriteFile("README.md", []byte("# test file"), 0777)
	if err != nil {
		t.Fatalf("could not write README.md: %v", err)
	}
	Git(t, "config", "user.email", "odog@git-commit-hook.ok")
	Git(t, "config", "user.name", "odog")
	Git(t, "add", "-A")
	Git(t, "commit", "-m", "initial commit")
	Git(t, "checkout", "-b", branchName)
}

// WriteConfigFile writes a test config file into the given folder
func WriteConfigFile(t *testing.T, dir string) {
	t.Helper()
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		t.Fatalf("could write config file, err in os.MkdirAll: %v", err)
	}
	cfg := config.Configuration{
		"test project": config.Project{
			Path: TestPathGitFolder,
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

// Git executes a git command
func Git(t *testing.T, args ...string) {
	t.Helper()
	o, err := exec.Command("git", args...).CombinedOutput()
	if err != nil {
		t.Fatalf("'git %v init' failed with error: %v - output: %s", args, err, string(o))
	}
}

// CleanupTestEnvironment removes all test folders
func CleanupTestEnvironment(t *testing.T) {
	err := os.RemoveAll(TestPath)
	if err != nil {
		t.Fatalf("Error cleaning up test environment.. Did not expect os.RemoveAll to return an error, but got: %v ", err)
	}
}

// PreapreTestEnvironment writes a test config file to the test folder
func PreapreTestEnvironment(t *testing.T) {
	WriteConfigFile(t, TestPath)
	err := os.Chdir(TestPath)
	if err != nil {
		t.Fatalf("Error preparing test environment.Did not expect os.Chdir to return an error, but got: %v ", err)
	}
}

// InitTestFolder cleans up an eventually existing folder, creates it again and changes working dir to that folder
func InitTestFolder(t *testing.T) {
	err := os.RemoveAll(TestPath)
	if err != nil {
		t.Fatalf("Did not expect os.RemoveAll to return an error, but got: %v ", err)
	}

	err = os.MkdirAll(TestPath, 0777)
	if err != nil {
		t.Fatalf("Did not expect os.MkdirAll to return an error, but got: %v ", err)
	}

	err = os.Chdir(TestPath)
	if err != nil {
		t.Fatalf("Did not expect os.Chdir to return an error, but got: %v ", err)
	}
}
