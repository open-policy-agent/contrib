VERSION := "0.1"
GIT_COMMIT := $(shell git rev-parse --short HEAD)
GO := go

CGO_ENABLED ?= 0
# Force modules on and to use the vendor directory.
GO := CGO_ENABLED=$(CGO_ENABLED) GO111MODULE=on GOFLAGS=-mod=vendor go

######################################################
#
# Development targets
#
######################################################

.PHONY: all
all: build check test

LDFLAGS := -X github.com/open-policy-agent/contrib/data_filter_elasticsearch/internal/version.Vcs=$(GIT_COMMIT) \
	-X github.com/open-policy-agent/contrib/data_filter_elasticsearch/internal/version.Version=$(VERSION)

.PHONY: build
build:
	$(GO) build -ldflags "$(LDFLAGS)" ./cmd/opa-es-filtering/...

.PHONY: check
check: check-fmt check-vet check-lint

.PHONY: check-fmt
check-fmt:
	./build/check-fmt.sh

.PHONY: check-vet
check-vet:
	./build/check-vet.sh

.PHONY: check-lint
check-lint:
	./build/check-lint.sh

.PHONY: test
test:
	$(GO) test ./...
