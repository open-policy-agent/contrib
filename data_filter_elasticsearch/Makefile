VERSION := "0.1"
GIT_COMMIT := $(shell git rev-parse --short HEAD)

all: build

LDFLAGS := -X github.com/open-policy-agent/contrib/data_filter_elasticsearch/internal/version.Vcs=$(GIT_COMMIT) \
	-X github.com/open-policy-agent/contrib/data_filter_elasticsearch/internal/version.Version=$(VERSION)

build:
	go build -ldflags "$(LDFLAGS)" ./cmd/opa-es-filtering/...
