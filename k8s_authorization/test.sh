#!/usr/bin/env bash

k() {
    kubectl --context kind-opa-authorizer "$@"
}

expect() {
    if [[ "$1" != "$2" ]]; then
        echo "Expected $1 == $2"
        exit 1
    fi
}

expect_ends_with() {
    if [[ "$1" != *"$2" ]]; then
        echo "Expected $1 == *$2"
        exit 1
    fi
}

echo "Waiting for OPA pod to come up"

exit_code=1
retries=0
opa_pod=""

while [[ "$exit_code" != 0 && "$retries" -lt 20 ]]
do
    sleep 5
    opa_pod=$(k --namespace opa get pods -l app=opa -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)
    exit_code="$?"
    ((retries++))
done

echo "OPA pod is up - awaiting condition=Ready"

# Wait for the OPA pod to become ready
k --namespace opa wait --for=condition=Ready --timeout=100s pods/"$opa_pod" > /dev/null

echo "OPA pod ready. Running tests."
echo "============================="

# Access to kube-system should be denied
result=$(k --namespace kube-system --as=someuser --as-group=system:authenticated get pods 2>&1)
expect "$?" 1
expect_ends_with "$result" "OPA: denied access to namespace kube-system"

# Access to opa namespace denied unless in devops group
result=$(k --namespace opa --as=someuser --as-group=system:authenticated get pods 2>&1)
expect "$?" 1
expect_ends_with "$result" "OPA: provided groups (system:authenticated) does not include all required groups: (devops, system:authenticated)"

# Access to opa namespace allowed if in devops group
result=$(k --namespace opa --as=someuser --as-group=system:authenticated --as-group=devops get pods 2>&1)
expect "$?" 0

echo "All tests successful!"
