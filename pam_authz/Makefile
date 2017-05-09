VERSION := 0.1-dev

.PHONY: all image up

all: image

image:
	docker build -t "pam-builder:$(VERSION)" -t "pam-builder:latest" \
			--force-rm -f docker/pam-builder.dockerfile .
	docker run "pam-builder" | \
		docker build -t "openpolicyagent/demo-pam:$(VERSION)" -t "openpolicyagent/demo-pam:latest" --force-rm -f run.dockerfile -

up:
	docker-compose -f docker/docker-compose.yaml up
