VERSION := 0.1-dev

.PHONY: all image up

all: image

image:
	cd docker; docker build -t openpolicyagent/demo-restful-api:latest -t openpolicyagent/demo-restful-api:$(VERSION) \
		-f Dockerfile .
up:
	docker-compose -f docker/docker-compose.yaml up
