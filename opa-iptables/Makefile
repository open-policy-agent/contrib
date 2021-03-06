VERSION := 0.0.1-dev
GO := go
GOOS := linux
DOCKER_TAG := 0.0.2-dev

BIN := opa-iptables
PKG := github.com/open-policy-agent/contrib/opa-iptables

BUILD_COMMIT := $(shell ./build/get-build-commit.sh)

LDFLAGS := "-X github.com/open-policy-agent/contrib/opa-iptables/pkg/version.Version=$(VERSION) \
	-X github.com/open-policy-agent/contrib/opa-iptables/pkg/version.Commit=$(BUILD_COMMIT) \
	-w"

GO111MODULE := on

.PHONY: all
all: test build push clean

.PHONY: version
version:
	@echo $(VERSION)

.PHONY: build
build: docker-build

.PHONY: docker-build
docker-build:
	@docker run --rm -v $$(pwd):/go/src/$(PKG) \
		-e GOOS=$(GOOS) \
		-e GO111MODULE=$(GO111MODULE) \
		-w /go/src/$(PKG) golang:1.13-rc-alpine  \
		$(GO) build -o $(BIN) -ldflags $(LDFLAGS)

	@docker build -t urvil38/opa-iptables:$(DOCKER_TAG) \
		-f Dockerfile \
		.

	@echo Successfully built $(BIN)

push:
	docker push urvil38/opa-iptables:$(DOCKER_TAG)

clean:
	rm -f $(BIN) -y

go-build:
	GOOS=$(GOOS) $(GO) build -o $(BIN) -ldflags $(LDFLAGS)

.PHONY: test
test: go-test

.PHONY: go-test
go-test:
	$(GO) test -v -count=1 ./...

