#!/usr/bin/env bash

k() {
    kubectl --context kind-opa-authorizer "$@"
}

pwd=$(pwd) envsubst < config/kind-conf.yaml | kind create cluster --name opa-authorizer --config -

k create namespace opa
k apply -k .
