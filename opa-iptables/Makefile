BIN := opa-iptables
PKG := github.com/contrib/opa-iptables
VERSION := 0.0.1-dev
ARCH := amd64
COMMIT := $(shell ./build/get-git-commit.sh)

.PHONY: all
all: build

.PHONY: build
build:
	/bin/sh -c "ARCH=$(ARCH) VERSION=$(VERSION) COMMIT=$(COMMIT) PKG=$(PKG) ./build/build.sh"