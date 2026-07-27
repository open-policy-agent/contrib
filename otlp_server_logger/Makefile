VERSION    := 0.1
REPOSITORY := openpolicyagent/otlp-server-logger

.PHONY: build test lint clean push

build:
	go build ./...
	go vet ./...
	go test ./...

test:
	go test ./...

lint:
	go vet ./...

push:
	docker buildx build --push --platform=linux/amd64,linux/arm64 \
		-f ./examples/Dockerfile \
		-t $(REPOSITORY):$(VERSION) \
		-t $(REPOSITORY):latest \
		.

clean:
	rm -f opa-otlp-server
