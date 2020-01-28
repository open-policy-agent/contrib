REPOSITORY := openpolicyagent/proxy_init
VERSION := v4

.PHONY: all
all: image

.PHONY: build
build: image

.PHONY: image
image:
	docker build -t $(REPOSITORY):latest -t $(REPOSITORY):$(VERSION) .

.PHONY: push
push: build
	docker push $(REPOSITORY):$(VERSION)
	docker push $(REPOSITORY):latest
