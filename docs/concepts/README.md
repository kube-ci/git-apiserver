# Concepts

Concepts help you learn about the different parts of the Git API server and the abstractions it uses.

- What is Git API server?
  - [Overview](/docs/concepts/what-is-git-apiserver/overview.md). Provides a conceptual introduction to Git API server, including the problems it solves and its high-level architecture.
- Custom Resource Definitions
  - [Repository](/docs/concepts/crds/repository.md). Introduces the concept of `Repository` for syncing a git repository in a Kubernetes native way.
  - [Branch](/docs/concepts/crds/branch.md). Introduces the concept of `Branch` to represent branches of git repositories.
  - [Tag](/docs/concepts/crds/tag.md). Introduce concept of `Tag` to represent tags of git repositories.
  - [PullRequest](/docs/concepts/crds/pull_request.md). Introduce concept of `PullRequest` to represent pull-requests of remote git repositories.
  - [GithubEvent](/docs/concepts/crds/github_event.md). Introduce concept of `GithubEvent` that represents Github events.