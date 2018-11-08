#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

GOPATH=$(go env GOPATH)
SRC=$GOPATH/src
BIN=$GOPATH/bin
ROOT=$GOPATH
REPO_ROOT=$GOPATH/src/github.com/kube-ci/git-apiserver

source "$REPO_ROOT/hack/libbuild/common/lib.sh"
source "$REPO_ROOT/hack/libbuild/common/kubeci_image.sh"

APPSCODE_ENV=${APPSCODE_ENV:-dev}
IMG=git-apiserver

DIST=$GOPATH/src/github.com/kube-ci/git-apiserver/dist
mkdir -p $DIST
if [ -f "$DIST/.tag" ]; then
  export $(cat $DIST/.tag | xargs)
fi

clean() {
  pushd $GOPATH/src/github.com/kube-ci/git-apiserver/hack/docker
  rm git-apiserver Dockerfile
  popd
}

build_binary() {
  pushd $GOPATH/src/github.com/kube-ci/git-apiserver
  ./hack/builddeps.sh
  ./hack/make.py build git-apiserver
  detect_tag $DIST/.tag
  popd
}

build_docker() {
  pushd $GOPATH/src/github.com/kube-ci/git-apiserver/hack/docker
  cp $DIST/git-apiserver/git-apiserver-alpine-amd64 git-apiserver
  chmod 755 git-apiserver

  cat >Dockerfile <<EOL
FROM alpine

RUN set -x \
  && apk add --update --no-cache ca-certificates

COPY git-apiserver /usr/bin/git-apiserver

USER nobody:nobody
ENTRYPOINT ["git-apiserver"]
EOL
  local cmd="docker build -t $DOCKER_REGISTRY/$IMG:$TAG ."
  echo $cmd; $cmd

  rm git-apiserver Dockerfile
  popd
}

build() {
  build_binary
  build_docker
}

docker_push() {
  if [ "$APPSCODE_ENV" = "prod" ]; then
    echo "Nothing to do in prod env. Are you trying to 'release' binaries to prod?"
    exit 0
  fi
  if [ "$TAG_STRATEGY" = "git_tag" ]; then
    echo "Are you trying to 'release' binaries to prod?"
    exit 1
  fi
  hub_canary
}

docker_release() {
  if [ "$APPSCODE_ENV" != "prod" ]; then
    echo "'release' only works in PROD env."
    exit 1
  fi
  if [ "$TAG_STRATEGY" != "git_tag" ]; then
    echo "'apply_tag' to release binaries and/or docker images."
    exit 1
  fi
  hub_up
}

source_repo $@
