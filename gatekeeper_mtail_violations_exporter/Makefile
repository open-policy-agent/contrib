VERSION := 1.0
REPOSITORY := openpolicyagent/gatekeeper_mtail_violatons_exporter

clean:
	rm -fr target

build: image

image:
	docker build -t $(REPOSITORY):$(VERSION) -t openpolicyagent/mtail_violatons_exporter:latest .

push: build
	docker push $(REPOSITORY):$(VERSION)
	docker push $(REPOSITORY):latest