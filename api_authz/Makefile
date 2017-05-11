REPOSITORY := openpolicyagent/demo-restful-api
VERSION := 0.1-dev

.PHONY: all
all: image

.PHONY: build
build: image

.PHONY: image
image:
	docker build -t $(REPOSITORY):latest \
		-t $(REPOSITORY):$(VERSION) \
		-f ./docker/Dockerfile \
		.

.PHONY: push
push: build
	docker push $(REPOSITORY):$(VERSION)
	docker push $(REPOSITORY):latest

.PHONY: up
up:
	docker-compose -f docker/docker-compose.yaml up
