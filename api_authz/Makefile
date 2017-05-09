VERSION := 0.1-dev

.PHONY: all image up

all: image

image:
	cd docker; docker build -t openpolicyagent/demo-http-api:latest -t openpolicyagent/demo-http-api:$(VERSION) \
		-f Dockerfile .
up:
	docker-compose -f docker/docker-compose.yaml up
