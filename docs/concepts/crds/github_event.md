# GithubEvent

## What is GithubEvent

A `GithubEvent` is a representation of Github [webhook event](https://developer.github.com/webhooks/#events) as a Kubernetes object with the help of [Aggregated API Servers](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/api-machinery/aggregated-api-servers.md). Currently, only [pull request events](https://developer.github.com/v3/activity/events/types/#pullrequestevent) are handled.

## GithubEvent structure

As with all other Kubernetes objects, a GithubEvent needs `apiVersion`, `kind`, and `metadata` fields. Here, we are going to describe some important sections of `GithubEvent` object.

### .action

The action that was performed.

### .repository

Describes the associated github repository.

### .sender

Describes the sender of this event.

### .pull_request

Describes the pull request itself. For more details see [here](https://developer.github.com/v3/activity/events/types/#pullrequestevent).

