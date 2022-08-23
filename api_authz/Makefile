REPOSITORY := openpolicyagent/demo-restful-api
VERSION := 0.3

.PHONY: all
all: image

.PHONY: build
build: image

.PHONY: image
image:
	docker build -t $(REPOSITORY):latest \
		-t $(REPOSITORY):$(VERSION) \
		./docker

.PHONY: push
push: build
	docker push $(REPOSITORY):$(VERSION)
	docker push $(REPOSITORY):latest

.PHONY: up
up:
	opa build docker/policy/example.rego
	docker-compose -f docker/docker-compose.yaml up
