#!/bin/sh

set -ex

GOARCH=${ARCH} go install -ldflags "-X ${PKG}/pkg/version.Version=${VERSION} -X ${PKG}/pkg/version.GitCommit=${COMMIT}" .