REPOSITORY := openpolicyagent/demo-pam
VERSION := 0.2-dev

.PHONY: image
all: image

.PHONY: build
build: image

.PHONY: image
image:
	docker build -t "pam-builder:$(VERSION)" \
		-t "pam-builder:latest" \
		--force-rm \
		-f ./docker/pam-builder.dockerfile \
		.
	docker run "pam-builder" | \
		docker build -t $(REPOSITORY):latest \
		-t $(REPOSITORY):$(VERSION) \
		--force-rm \
		-f ./run.dockerfile \
		-

.PHONY: push
push: build
	docker push $(REPOSITORY):$(VERSION)
	docker push $(REPOSITORY):latest

.PHONY: up
up:
	docker-compose -f docker/docker-compose.yaml up
