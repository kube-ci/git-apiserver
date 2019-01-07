[![Go Report Card](https://goreportcard.com/badge/github.com/kube-ci/git-apiserver)](https://goreportcard.com/report/github.com/kube-ci/git-apiserver)
[![Build Status](https://travis-ci.org/kube-ci/git-apiserver.svg?branch=master)](https://travis-ci.org/kube-ci/git-apiserver)
[![codecov](https://codecov.io/gh/kube-ci/git-apiserver/branch/master/graph/badge.svg)](https://codecov.io/gh/kube-ci/git-apiserver)
[![Docker Pulls](https://img.shields.io/docker/pulls/kubeci/git-apiserver.svg)](https://hub.docker.com/r/kubeci/git-apiserver/)
[![Slack](https://slack.appscode.com/badge.svg)](https://slack.appscode.com)
[![Twitter](https://img.shields.io/twitter/follow/thekubeci.svg?style=social&logo=twitter&label=Follow)](https://twitter.com/intent/follow?screen_name=TheKubeCi)

[![Throughput Graph](https://graphs.waffle.io/kube-ci/project/throughput.svg)](https://waffle.io/kube-ci/project/metrics/throughput)

# Git API Server

Git API server by AppsCode is a Kubernetes operator for syncing Git repositories as Kubernetes resources.

## Features

- Sync branches and tags of a git repository.
- Sync pull-requests using webhook events.
- Configure credentials for syncing private repositories.

## Supported Versions

Please pick a version of Git API server that matches your Kubernetes installation.

| Git API server Version                                                                      | Docs                                                            | Kubernetes Version |
|------------------------------------------------------------------------------------|-----------------------------------------------------------------|--------------------|
| [0.1.0](https://github.com/kube-ci/git-apiserver/releases/tag/0.1.0) (uses CRD) | [User Guide](https://kube.ci/products/git-apiserver/0.1.0)    | 1.9.x+             |

## Installation

To install Git API server, please follow the guide [here](https://kube.ci/products/git-apiserver/0.1.0/setup/install).

## Using Git API Server

Want to learn how to use Git API server? Please start [here](https://kube.ci/products/git-apiserver/0.1.0).

## Git API Server API Clients

You can use Git API server api clients to programmatically access its objects. Here are the supported clients:

- Go: [https://github.com/kube-ci/git-apiserver](/client/clientset/versioned)

## Contribution guidelines

Want to help improve Git API server? Please start [here](https://kube.ci/products/git-apiserver/0.1.0/welcome/contributing).

---

**Git API server binaries collects anonymous usage statistics to help us learn how the software is being used and how we can improve it. To disable stats collection, run the operator with the flag** `--enable-analytics=false`.

---

## Support

We use Slack for public discussions. To chit chat with us or the rest of the community, join us in the [AppsCode Slack team](https://appscode.slack.com/messages/C8NCX6N23/details/) channel `#kubeci`. To sign up, use our [Slack inviter](https://slack.appscode.com/).

If you have found a bug with Git API server or want to request for new features, please [file an issue](https://github.com/kube-ci/project/issues/new).
