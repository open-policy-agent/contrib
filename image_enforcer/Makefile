BIN 		:= clair-layer-sync
PKG 		:= github.com/open-policy-agent/contrib/image_enforcer
REPOSITORY 	:= openpolicyagent/$(BIN)
VERSION 	:= 0.1-dev

build:
	docker run --rm -v $$(pwd):/go/src/$(PKG) \
		-w /go/src/$(PKG) golang:1.8-alpine \
		go build -o $(BIN)_linux_amd64 .
	docker build -t $(REPOSITORY):$(VERSION) \
		-t $(REPOSITORY):latest \
		-f Dockerfile.run \
		.
	@echo Successfully built $(REPOSITORY):$(VERSION) $(REPOSITORY):latest

push:
	docker push $(REPOSITORY)

clean:
	rm -f $(BIN)_*_*
