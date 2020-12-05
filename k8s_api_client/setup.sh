#!/usr/bin/env bash

k() {
    kubectl --context kind-opa-kubernetes-api-client "$@"
}

kind create cluster --name opa-kubernetes-api-client --config kind-config.yaml

# Install nginx ingress controller and wait for it to become ready
k apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/provider/kind/deploy.yaml

# Hackish attempt to avoid "error: no matching resources found" in next call
sleep 5

k wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector="app.kubernetes.io/component=controller" \
  --timeout=180s

k apply -k .
