VERSION := $(shell ./scripts/get-version.sh)

clean:
	rm -fr target

build: jar image

image:
	sed s/VERSION/$(VERSION)/g Dockerfile.in > Dockerfile
	docker build -t openpolicyagent/demo-kafka:$(VERSION) .

jar:
	./scripts/build.sh
