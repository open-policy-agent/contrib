# opa-kubernetes-api-client

Rego functions to query the Kubernetes API directly from OPA.

```rego
package authz

import data.kubernetes.api.client

mutating_actions := {"update", "patch", "delete"}

# Deny mutating action unless user is in group owning the resource
deny[reason] {
    mutating_actions[input.spec.resourceAttributes.verb]

    cluster_resource := client.query_name_ns(
        input.spec.resourceAttributes.resource,
        input.spec.resourceAttributes.name,
        input.spec.resourceAttributes.namespace,
    )
    resource_owner := cluster_resource.metadata.labels.owner
    user_groups := {group | group := input.spec.user.group[_]}

    not user_groups[resource_owner]

    reason := sprintf("User %v not in group %v, has %v", [input.spec.user.name, resource_owner, user_groups])
}
```

## Motivation

Many policies interacting with Kubernetes depends on knowing the current state of the cluster. Examples of this include checking for hostname or paths already in use by ingress resources, the number of workloads running on a node, or the total amount of memory allocated to an application. Having this data available for decisions significantly increases the number of possible use cases for policy enforcement, whether it be for authorization or admission control.

Systems that integrate OPA and Kubernetes commonly provide a [cached](https://github.com/open-policy-agent/kube-mgmt#caching) or [replicated](https://github.com/open-policy-agent/gatekeeper#replicating-data) view of the Kubernetes API in order to enable these type of policies. In some cases, replicating cluster state into OPA is overkill and it's more efficient to simply query for the required resource(s) from within the policy. In other cases, resources need to be queried at policy evaluation time to obtain a state as close to the time of evaluation as possible. Note, because (by design) Kubernetes is eventually consistent, it's still possible for the API call to return stale results. In some cases you might even want to do both - consult the cache first for performance, and perform an API lookup only if an object is uncached or the cache entry is suspected to be stale.

## Setup

In order to communicate with the kubernetes API you'll need a service account and a service account token. Since OPA can't read the token from disk there's no need to mount it there. Instead, we'll put it in a secret which we can mount as an environment variable.

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: opa-kubernetes-api-client
automountServiceAccountToken: false
---
apiVersion: v1
kind: Secret
metadata:
  name: opa-kubernetes-api-client
  annotations:
    kubernetes.io/service-account.name: opa-kubernetes-api-client
type: kubernetes.io/service-account-token
```

This will have the token controller automatically create a token for us and assign it to the `opa-kubernetes-api-client` service account. We'll now need to make the token visible to our policies - for this we'll patch our OPA container definitions to expose our token as an environment variable.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: opa
spec:
  template:
    spec:
      containers:
      - name: opa
        env:
        - name: KUBERNETES_API_TOKEN
          valueFrom:
            secretKeyRef:
              name: opa-kubernetes-api-client
              key: token
```

We may now retrieve the service account token from inside our rego policies by calling `opa.runtime().env.KUBERNETES_API_TOKEN`.

### RBAC

In order to actually read any data, you'll need to allow your service account to do so. For this you'll need to create an RBAC role and role binding. In the example below we'll create a role that has cluster level read access to some common resource types, but you'll naturally want to tweak this for your needs:

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: opa-kubernetes-api-client
rules:
- apiGroups: [""]
  resources: ["services", "pods", "configmaps"]
  verbs: ["get", "list"]
- apiGroups: ["apps"]
  resources: ["daemonsets", "deployments", "statefulsets", "replicasets"]
  verbs: ["get", "list"]
```

Let's bind the role to our service account:

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: opa-kubernetes-api-client-read
subjects:
- kind: ServiceAccount
  name: opa-kubernetes-api-client
  namespace: default
roleRef:
  kind: ClusterRole
  name: opa-kubernetes-api-client
  apiGroup: rbac.authorization.k8s.io
```

All done! You should now be able to call the functions included in this library to query the kubernetes API for resources directly.

## Available functions

```rego
# Query for given resource/name in provided namespace
# Example: query_ns("deployments", "my-app", "default")
query_name_ns(resource, name, namespace)

# Query for given resource type using label selectors
# https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#api
# Example: query_label_selector_ns("deployments", {"app": "opa-kubernetes-api-client"}, "default")
query_label_selector_ns(resource, selector, namespace)

# Query for all resources of type resource in all namespaces
# Example: query_all("deployments")
query_all(resource)
```

## Local development and testing

To develop and/or test the client locally, the following applications need to be installed:

* [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
* [curl](https://curl.haxx.se/download.html)

Run the `setup.sh` script to create a local kubernetes cluster using [kind](https://kind.sigs.k8s.io/). This will deploy OPA and an ingress controller to let you access it as if running normally (i.e. on `localhost:8181`). You can then run queries using the [OPA query API](https://www.openpolicyagent.org/docs/latest/rest-api/#query-api) in order to test the functions from the policy.

`POST http://localhost:8181/v1/query`
```json
{
    "query": "x := data.kubernetes.api.client.query_name_ns(\"deployments\",\"opa-kubernetes-api-client\", \"default\").body"
}
```

`POST http://localhost:8181/v1/query`
```json
{
    "query": "x := data.kubernetes.api.client.query_all(\"deployments\").body"
}
```

`POST http://localhost:8181/v1/query`
```json
{
    "query": "x := data.kubernetes.api.client.query_label_selector_ns(\"deployments\", {\"app\":\"opa-kubernetes-api-client\"}, \"default\").body"
}
```
Use the `test.sh` script to run these queries against the test cluster and verify the expected results.