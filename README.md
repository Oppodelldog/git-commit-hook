# git-commit-hook

[![Go Report Card](https://goreportcard.com/badge/github.com/Oppodelldog/git-commit-hook)](https://goreportcard.com/report/github.com/Oppodelldog/git-commit-hook) [![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://raw.githubusercontent.com/Oppodelldog/git-commit-hook/master/LICENSE) [![Linux build](http://nulldog.de:12080/api/badges/Oppodelldog/git-commit-hook/status.svg)](http://nulldog.de:12080/Oppodelldog/git-commit-hook) [![Windows build](https://ci.appveyor.com/api/projects/status/qpe2889fbk1bw7lf/branch/master?svg=true)](https://ci.appveyor.com/project/Oppodelldog/git-commit-hook/branch/master) [![Coverage Status](https://coveralls.io/repos/github/Oppodelldog/git-commit-hook/badge.svg?branch=master)](https://coveralls.io/github/Oppodelldog/git-commit-hook?branch=master)

**Problem:** I forget to add ticket numbers to commit messages or I take the wrong ticket ID.

**Solution:** a custom git-hook that prepends the current branch-name to every commit message.

The implementation and rules are tightly coupled to **git flow** used with **jira**.

The intention is on the one hand to make life easier when working on features.
On the other hand the intention must be to not allow commits to non-fetaure branches without having a valid ticket reference.

### So here are the rules

* At least there must always be a non empty commit message.

* When you are committing to a **feature** branch

        the hook will preprend the branch name to the commit message

* When you are committing to a non feature branch, lets say **develop** for a quick fix, **release/v0.1.1** for a release fix or a **hotfix**

        you have to enter a valid feature reference.
