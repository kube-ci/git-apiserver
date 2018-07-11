#!/usr/bin/env bash

pushd $GOPATH/src/kube.ci/git-apiserver/hack/gendocs
go run main.go
popd
