VERSION := 0.1-dev

.PHONY: all image up

all: image

image:
	docker build -t "pam-builder:$(VERSION)" -t "pam-builder:latest" \
			--force-rm -f docker/pam-builder.dockerfile .
	docker run "pam-builder" | \
		docker build -t "openpolicyagent/ssh-webapp:$(VERSION)" -t "openpolicyagent/demo-pam:latest" --force-rm -f run.dockerfile -
	# docker run "pam-builder" | \
	# 	docker build -t "openpolicyagent/ssh-backend:$(VERSION)" -t "openpolicyagent/ssh-backend:latest" --force-rm --build-arg identifier=Backend -f run.dockerfile -

up:
	docker-compose -f docker/docker-compose.yaml up
