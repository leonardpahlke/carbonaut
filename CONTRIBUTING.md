# Contributing to Carbonaut

This file describes how you can contribute to the project. We follow a common git workflow to make changes that drives with opening pull requests from forked repositories.

## Communication Channels

*tbd*

## Meetings

*tbd*

## Project Governance

Project governance is documented [here](./GOVERNANCE.md).

## GitHub Workflow

The recommended workflow is to fork this repository and open pull requests from your fork.

![contribution-workflow](/img/carbonaut-contributing-workflow-v3.png)

### 1. Fork, clone & configure Carbonaut upstream

- Click on the _Fork_ [button on GitHub]( https://docs.github.com/en/get-started/quickstart/fork-a-repo)
- Clone your fork
- Add the upstream repository as a new remote

```sh
# Clone repository
git clone https://github.com/$YOUR_GITHUB_USER/$REPOSITORY.git

# Add upstream origin
git remote add upstream git@github.com:carbonaut-cloud/$REPOSITORY.git
```

### 2. Create a pull request

```sh
# Create a new feature branch
git checkout -b my_feature_branch

# Make changes to your branch
# ...

# Commit changes - remember to sign!
git commit -s

# Push your new feature branch
git push --set-upstream origin my_feature_branch

# Create a new pull request from https://github.com/carbonaut-cloud/$REPOSITORY
```

### 3. Update your pull request with latest changes

```sh
# Checkout main branch
git checkout main

# Update your fork's main branch from upstream
git pull upstream main

# Checkout your feature branch
git checkout my_feature_branch

# Rebase your feature branch changes on top of the updated main branch
git rebase main

# Update your pull request with latest changes
git push -f my_feature_branch
```

## Scope of pull requests

We prefer small incremental changes that can be reviewed and merged quickly.
It's OK if it takes multiple pull requests to close an issue.

The idea is that each improvement should land in Carbonaut's main branch within a
few hours.  The sooner we can get multiple people looking at and agreeing on a
specific change, the quicker we will have it out in a release.  The quicker we
can get these small improvementes in a Carbonaut release, the quicker we can get
feedback from our users and find out what doesn't work, or what we have missed.

The added benefit is that this will force everyone to think about handling
partially implemented features & non-breaking changes.

## Commits

### DCO

Contributions to this project must be accompanied by a Developer Certificate of
Origin (DCO).

All commit messages must contain the Signed-off-by line with an email address that matches the commit author. When commiting, use the `--signoff` flag:

```sh
git commit -s
```

You can find more about this [here]( https://docs.github.com/en/authentication/managing-commit-signature-verification/signing-commits).

The Signed-off-by line must match the **author's real name**, otherwise the PR will be rejected.

### Commit messages

[How to Write a Git Commit Message](https://chris.beams.io/posts/git-commit/)

Guidelines:

- **Group Commits:** Each commit should represent a meaningful change (e.g. implement feature X, fix bug Y, ...).
  - For instance, a PR should not look like _1) Add Feature X 2) Fix Typo 3) Changes to features X 5) Bugfix for feature X 6) Fix Linter 7) ...
  - Instead, these commits should be squashed together into a single "Add Feature" commit.
- Each commit should work on its own: it must compile, pass the linter and so on.
  - This makes life much easier when using `git log`, `git blame`, `git bisect`, etc.
  - For instance, when doing a `git blame` on a file to figure out why a change
  was introduced, it's pretty meaningless to see a _Fix linter_ commit message.
  "Add Feature X" is much more meaningful.
- Use `git rebase -i main` to group commits together and rewrite their commit message
- To add changes to the previous commit, use `git commit --amend -s`. This will
  change the last commit (amend) instead of creating a new commit.


[^1]: See [https://www.conventionalcommits.org](https://www.conventionalcommits.org)

## Docs

### Use relative links to markdown files

Link to markdown files (`[link](../foo.md)`) instead of relative URLs
(`[link](/foo)`).

The docs compiler will replace file links with relative URLs automatically.

This is to avoid broken links. If a file gets renamed, the compiler will
catch broken links and throw an error. Relative URLs get broken unnoticed.
