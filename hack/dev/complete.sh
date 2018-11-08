#!/usr/bin/env bash
set -xe

REPO=${GOPATH}/src/github.com/kube-ci/git-apiserver
pushd ${REPO}

export APPSCODE_ENV=dev

# codegen
./hack/codegen.sh

# make.py
./hack/make.py

# build docker
./hack/docker/setup.sh

# load to minikube
minducker kubeci/git-apiserver:initial

# deploy
./hack/deploy/install.sh
