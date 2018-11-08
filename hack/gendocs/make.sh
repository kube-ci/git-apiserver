#!/usr/bin/env bash

pushd $GOPATH/src/github.com/kube-ci/git-apiserver/hack/gendocs
go run main.go
popd
