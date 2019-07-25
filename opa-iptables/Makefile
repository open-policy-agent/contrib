VERSION := 1.0.0
GO := go
GOOS := linux

BIN := opa-iptables
PKG := github.com/open-policy-agent/contrib/opa-iptables

BUILD_COMMIT := $(shell ./build/get-build-commit.sh)

LDFLAGS := "-X github.com/open-policy-agent/contrib/opa-iptables/pkg/version.Version=$(VERSION) \
	-X github.com/open-policy-agent/contrib/opa-iptables/pkg/version.Commit=$(BUILD_COMMIT) \
	-w"

GO111MODULE := on

.PHONY: all
all: test build

.PHONY: version
version:
	@echo $(VERSION)

.PHONY: build
build: docker-build

.PHONY: docker-build
docker-build:
	@docker run -it --rm -v $$(pwd):/go/src/$(PKG) \
		-e GOOS=$(GOOS) \
		-e GO111MODULE=$(GO111MODULE) \
		-w /go/src/$(PKG) urvil38/golang:1.12  \
		$(GO) build -o $(BIN) -ldflags $(LDFLAGS)

	@echo Successfully built $(BIN)

.PHONY: test
test: go-test

.PHONY: go-test
go-test:
	$(GO) test -v -count=1 ./...

