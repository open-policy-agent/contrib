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

LDFLAGS := -X github.com/open-policy-agent/contrib/data_filter_mongodb/internal/version.Vcs=$(GIT_COMMIT) \
	-X github.com/open-policy-agent/contrib/data_filter_mongodb/internal/version.Version=$(VERSION)

.PHONY: build-binary
build-binary:
	go build -o bin/opa-mongo github.com/open-policy-agent/contrib/data_filter_mongodb

.PHONY: vet
vet:
	go vet ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	$(GO) test ./...

.PHONY: build
build:
	docker build -t $(REPOSITORY):latest -t $(REPOSITORY):$(VERSION) .

.PHONY: push
push: build
	docker push $(REPOSITORY):$(VERSION)
	docker push $(REPOSITORY):latest

all: build-binary vet fmt lint test build push


