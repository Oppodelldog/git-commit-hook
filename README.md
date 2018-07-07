[![Go Report Card](https://goreportcard.com/badge/github.com/Oppodelldog/git-commit-hook)](https://goreportcard.com/report/github.com/Oppodelldog/git-commit-hook) [![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://raw.githubusercontent.com/Oppodelldog/git-commit-hook/master/LICENSE) [![Linux build](http://nulldog.de:12080/api/badges/Oppodelldog/git-commit-hook/status.svg)](http://nulldog.de:12080/Oppodelldog/git-commit-hook)

# git-commit-hook
> configureable git commit hook

### 1. Install

Downlod the binary, put it into a folder of your $PATH (for example. /usr/local/bin).

Ensure your user has execution permissions on that file.

### 2. Configure
There are several places you can put the configuration.

Create a config file named ```git-commit-hook.yaml```.

Create the configuration below your user folder:

**/home/*user*/.config/git-commit-hook**


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
        "^.*(#\\d*)+.*$" : must have a ticket reference (eg. #267)
 ```

Watch out the test [fixture](config/test-data.yaml) for full feature sample

### 3. Activate
Use the subcommand ```install``` to activate the commit-message-hook in your repository.


## Sub-Commands
There are some useful subcommands which ease the use of the commit-message hook.

### git-commit-hook install
Installs the commit-hook in the configured repositories.

You need to specify either **-p** or **-a**.

* **-p** to install in the given repository (eg. **-p "project xyz"**)
* **-a** to install in all configured repositories

If there's already a commit-message-hook installed, you can overwrite by adding ```-f```.

### git-commit-hook uninstall
Uninstalls the commit-hook from the configured repositories.

You need to specify either **-p** or **-a**.

* **-p** to uninstall from the given repository (eg. **-p "project xyz"**)
* **-a** to uninstall from all configured repositories

### git-commit-hook diag
Gives an overview of the configuration and the installed commit hooks

### git-commit-hook test

The test command is useful to test configuration and simulate a commit-situation.

You may input the following parameters:
* **-p** project name
* **-b** branch name
* **m** commit message

**Sample:**

```shell
git-commit-hook test -m "short commit message" -b master -p testrepo
```

in this case, validation will fail, since it's required to give a ticket/issue reference in the commit message.

```shell
testing configuration '/home/nils/.config/git-commit-hook/git-commit-hook.yaml':
project        : testrepo
branch         : master
commit message : short commit message

validation error for branch 'master'
at least expected one of the following to match
 - must have a ticket reference (eg. #267)
```

If you pass a valid commit message
```shell
git-commit-hook test -m "fix #121" -b master -p testrepo
```

shows not validation errors, but shows the finally rendered message.
```shell
testing configuration '/home/nils/.config/git-commit-hook/git-commit-hook.yaml':
project        : testrepo
branch         : master
commit message : fix #121

would generate the following commit message:
fix #121
```