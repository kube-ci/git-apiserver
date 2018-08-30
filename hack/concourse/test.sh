#!/usr/bin/env bash

set -eoux pipefail

ORG_NAME=kube-ci
REPO_NAME=git-apiserver
APP_LABEL=kube-ci #required for `kubectl describe deploy -n kube-system -l app=$APP_LABEL`

export APPSCODE_ENV=dev
export DOCKER_REGISTRY=appscodeci

# get concourse-common
pushd $REPO_NAME
git status # required, otherwise you'll get error `Working tree has modifications.  Cannot add.`. why?
git subtree pull --prefix hack/libbuild https://github.com/appscodelabs/libbuild.git master --squash -m 'concourse'
popd

source $REPO_NAME/hack/libbuild/concourse/init.sh

pushd $GOPATH/src/github.com/$ORG_NAME/$REPO_NAME

# install dependencies
./hack/builddeps.sh

# run tests
./hack/make.py test e2e

popd
