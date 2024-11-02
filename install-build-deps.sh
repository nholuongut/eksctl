#!/bin/sh -eux

# Make sure to bump the version of EKSCTL_DEPENDENCIES_IMAGE if you make any changes here

if [ -z "${GOBIN+x}" ]; then
 GOBIN="$(go env GOPATH)/bin"
fi

if [ "$(uname)" = "Darwin" ] ; then
  OSARCH="darwin-amd64"
else
  OSARCH="linux-amd64"
fi

env CGO_ENABLED=1 go install -tags extended github.com/gohugoio/hugo

go install github.com/jteeuwen/go-bindata/go-bindata
go install github.com/vektra/mockery/cmd/mockery
go install github.com/nholuongut/github-release
go install golang.org/x/tools/cmd/stringer

# TODO: metalinter is archived, we should switch to github.com/golangci/golangci-lint
# Install metalinter
# Managing all linters that gometalinter uses with dep is going to take
# a lot of work, so we install all of those from the release tarball
METALINTER_VERSION="3.0.0"
curl --silent --location "https://github.com/alecthomas/gometalinter/releases/download/v${METALINTER_VERSION}/gometalinter-${METALINTER_VERSION}-${OSARCH}.tar.gz" \
  | tar -x -z -C "${GOBIN}" --strip-components 1
