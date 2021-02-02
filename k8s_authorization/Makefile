.PHONY: deploy
deploy:
	./setup.sh

.PHONY: test
test:
	./test.sh

.PHONY: cleanup
cleanup:
	kind delete cluster --name opa-authorizer

.PHONY: build
build: deploy test cleanup
