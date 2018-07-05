[![Go Report Card](https://goreportcard.com/badge/github.com/Oppodelldog/git-commit-hook)](https://goreportcard.com/report/github.com/Oppodelldog/git-commit-hook) [![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://raw.githubusercontent.com/Oppodelldog/git-commit-hook/master/LICENSE) [![Linux build](http://nulldog.de:12080/api/badges/Oppodelldog/git-commit-hook/status.svg)](http://nulldog.de:12080/Oppodelldog/git-commit-hook) [![Windows build](https://ci.appveyor.com/api/projects/status/qpe2889fbk1bw7lf/branch/master?svg=true)](https://ci.appveyor.com/project/Oppodelldog/git-commit-hook/branch/master) [![Coverage Status](https://coveralls.io/repos/github/Oppodelldog/git-commit-hook/badge.svg?branch=master)](https://coveralls.io/github/Oppodelldog/git-commit-hook?branch=master)

# git-commit-hook
> configureable git commit hook


**Customize commit messages dynamically with templating**

**Validate commit message**

### 1. install
#### Download
downlod the binary and ensure your user has execution permissions
#### Install
Copy or link the binary into your git repositories ```hooks``` folder
(Project/.git/hooks), rename it to ```commit-msg```

Or create a symlink:
```ln -sf ~/Downloads/git-commit-hook ~/MyProject/.git/hooks/commit-msg```

### 2. Configure
```yaml
 "project xyz":
   # path to the git repository
   path: "/home/nils/projects/xyz/.git"

   # define types of branch and a pattern to identify the type for the current branch you are working on
   branch:
      master: "^(origin\\/)*master"

   # define a commit message template per branch type, or as here for all (*) branch types
   template:
     "*": "{.BranchName}: {.CommitMessage}"

   # define validation rules per branch type
   validation:
      master:
        "(?m)(?:\\s|^|/)(([A-Z](_)*)+-[0-9]+)([\\s,;:!.-]|$)" : "valid ticket ID"
 ```
 There are several places the configuration will be searched at, but one thing is for sure, the config file
 must be named ```git-commit-hook.yaml```.

 Here are some places the configuration will be searched at:
 * **~/.config/git-commit-hook**
 * inside the **git repository** you commit in (also in subfolders **.git**, **.git/hooks**)

Watch out the test [fixture](config/test-data.yaml) for full feature sample

### 3. Commit
You can do it on your own, I know that!!

---

### Configuration check
to check if the git hook is installed and configured correctly, just call the command from
your repository like this:
```.git/hooks/commit-msg```

The output will be something like this:

    git-commit-hook - parsed configuration


    branch types:
         master : ^(origin\/)*master

    branch type templates:
         master : {{.CommitMessage}} - whatever !!! {{.BranchName}}

    branch type validation:
         master :
             (?m)^.*(#\d*)+.*$ : must have a ticket reference (eg. #267)
