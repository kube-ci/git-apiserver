## Development Guide

This document is intended to be the canonical source of truth for things like supported toolchain versions for building Git API server.
If you find a requirement that this doc does not capture, please submit an issue on github.

This document is intended to be relative to the branch in which it is found. It is guaranteed that requirements will change over time
for the development branch, but release branches of Git API server should not change.

### Build Git API server

Some of the Git API server development helper scripts rely on a fairly up-to-date GNU tools environment, so most recent Linux distros should
work just fine out-of-the-box.

#### Setup GO

Git API server is written in Google's GO programming language. Currently, Git API server is developed and tested on **go 1.10**. If you haven't set up a GO
development environment, please follow [these instructions](https://golang.org/doc/code.html) to install GO.

#### Download Source

```console
$ cd $(go env GOPATH)/src/github.com/kube-ci
$ git clone https://github.com/kube-ci/git-apiserver.git
$ cd git-apiserver
```

#### Install Dev tools

To install various dev tools for Git API server, run the following command:
```console
$ ./hack/builddeps.sh
```

#### Build Binary

```
$ ./hack/make.py
$ git-apiserver version
```

#### Run Binary Locally

```console
$ git-apiserver run \
  --secure-port=8443 \
  --kubeconfig="$HOME/.kube/config" \
  --authorization-kubeconfig="$HOME/.kube/config" \
  --authentication-kubeconfig="$HOME/.kube/config" \
  --authentication-skip-lookup
```

#### Dependency management

Git API server uses [Glide](https://github.com/Masterminds/glide) to manage dependencies. Dependencies are already checked in the `vendor` folder.
If you want to update/add dependencies, run:
```console
$ glide slow
```

#### Build Docker images

To build and push your custom Docker image, follow the steps below. To release a new version of Git API server, please follow the [release guide](/docs/setup/developer-guide/release.md).

```console
# Build Docker image
$ ./hack/docker/setup.sh; ./hack/docker/setup.sh push

# Add docker tag for your repository
$ docker tag kubeci/git-apiserver:<tag> <image>:<tag>

# Push Image
$ docker push <image>:<tag>
```

#### Generate CLI Reference Docs

```console
$ ./hack/gendocs/make.sh
```

### Testing Git API server

#### Unit tests

```console
$ ./hack/make.py test unit
```

#### Run e2e tests

Git API server uses [Ginkgo](http://onsi.github.io/ginkgo/) to run e2e tests.
```console
$ ./hack/make.py test e2e
```

To run e2e tests against remote backends, you need to set cloud provider credentials in `./hack/config/.env`. You can see an example file in `./hack/config/.env.example`.
