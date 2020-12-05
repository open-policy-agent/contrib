#!/usr/bin/env bash

k() {
    kubectl --context kind-opa-kubernetes-api-client "$@"
}

k wait --for=condition=ready pod \
  --selector="app=opa-kubernetes-api-client" \
  --timeout=30s

query='{"query": "x := data.kubernetes.api.client.query_name_ns(\"deployments\",\"opa-kubernetes-api-client\", \"default\").body"}'
result=$(curl --silent --data "$query" http://localhost:8181/v1/query | jq -r .result[].x.metadata.name)
if [[ "$result" != "opa-kubernetes-api-client" ]]; then
    echo "Expected metadata.name 'opa-kubernetes-api-client' for query $query"
    exit 1
fi

query='{"query": "x := data.kubernetes.api.client.query_all(\"deployments\").body"}'
result=$(curl --silent --data "$query" http://localhost:8181/v1/query | jq -r '.result[].x.items[0]'.metadata.name)
if [[ "$result" != "opa-kubernetes-api-client" ]]; then
    echo "Expected metadata.name 'opa-kubernetes-api-client' for query $query"
    exit 1
fi

query='{"query": "x := data.kubernetes.api.client.query_label_selector_ns(\"deployments\", {\"app\":\"opa-kubernetes-api-client\"}, \"default\").body"}'
result=$(curl --silent --data "$query" http://localhost:8181/v1/query | jq -r '.result[].x.items[0]'.metadata.name)
if [[ "$result" != "opa-kubernetes-api-client" ]]; then
    echo "Expected metadata.name 'opa-kubernetes-api-client' for query $query"
    exit 1
fi

echo "All tests successfully executed"