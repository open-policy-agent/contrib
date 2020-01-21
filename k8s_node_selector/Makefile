KIND_CLUSTER_NAME ?= "opa-rego-teammapping"

WAIT_FOR_KIND_READY = '{range .items[*]}{@.metadata.name}:{range @.status.conditions[*]}{@.type}={@.status};{end}{end}'
WAIT_FOR_OPA_READY = '{range .items[*]}{@.metadata.name}:{range @.status.conditions[*]}{@.type}={@.status};{end}{end}'

.PHONY: all build test push kind-start kind-deploy e2e deocker-e2e integration integration-update-rego logs
all: e2e

build: test

test:
	opa test ./*.rego -v --explain full

push: 
	@echo "No assets to push as deployed to k8s via config map"

kind-start:
ifeq (1, $(shell kind get clusters | grep ${KIND_CLUSTER_NAME} | wc -l))
	@echo "Cluster already exists - deleting it to start from clean cluster" 
	kind delete cluster --name ${KIND_CLUSTER_NAME}
endif

	@echo "Creating Cluster"	
	kind create cluster --name ${KIND_CLUSTER_NAME} --image=kindest/node:v1.16.2
	# Wait for cluster to be ready
	until kubectl get nodes -o jsonpath="${WAIT_FOR_KIND_READY}" 2>&1 | grep -q "Ready=True"; do sleep 5; echo "--------> waiting for cluster node to be available"; done

kind-deploy: kind-start
	./integration/deploy/deploy.sh
	# Wait for OPA to be ready
	until kubectl -n opa -lapp=opa get pods -o jsonpath="${WAIT_FOR_OPA_READY}" 2>&1 | grep -q "Ready=True"; do sleep 5;echo "--------> waiting for operator to be available"; kubectl get pods -n opa; done


e2e: kind-start kind-deploy integration

docker-e2e:
	docker build -t opa-rego-teammapping-testcontainer -f ./.devcontainer/Dockerfile .
	docker run -v ${PWD}:/src --workdir /src -v /var/run/docker.sock:/var/run/docker.sock --network=host opa-rego-teammapping-testcontainer make e2e

integration: test
	# Run integration tests
	python3 ./integration/test/int_test.py

# Normal update of the rego policy
update-rego: test
	# Delete the existing rego config map
	kubectl delete configmap opa-default-system-main -n opa || true
	# Upload new one
	kubectl create configmap opa-default-system-main --from-file ./main.rego -n opa

# Forced update of the rego policy. Useful if the update has caused a crash loop or other
# issue which prevent normal update
force-update-rego: test
	# Delete the existing rego config map
	kubectl delete configmap opa-default-system-main -n opa || true
	# Upload new one
	kubectl create configmap opa-default-system-main --from-file ./main.rego -n opa
	# Delete OPA Agent pod to force policy reload
	kubectl get pods -n opa | grep opa- | awk '{print $$1}' | xargs kubectl delete -n opa pod

logs:
	kubectl get pods -n opa | grep opa- | awk '{print $$1}' | xargs -I {} kubectl logs {} opa

logs-sidecar:
	kubectl get pods -n opa | grep opa- | awk '{print $$1}' | xargs -I {} kubectl logs {} kube-mgmt