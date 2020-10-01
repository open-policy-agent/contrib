POLICIES    := policies
OPA_IMAGE   := openpolicyagent/opa
OPA_VERSION := 0.23.2

test:
	docker run --rm -v $$(PWD)/$(POLICIES):/tmp $(OPA_IMAGE):$(OPA_VERSION) test /tmp -v
