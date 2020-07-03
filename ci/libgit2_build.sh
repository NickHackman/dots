#!/usr/bin/env bash
go get -d github.com/libgit2/git2go 
cd "${GOPATH}"/src/github.com/libgit2/git2go || exit 1
git checkout next
git submodule update --init # get libgit2
make install

