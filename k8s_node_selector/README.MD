# K8s Node Selector

Cloud providers support the concept of `agentpools` or `nodepools` in a Kubernetes ([Azure AKS](https://docs.microsoft.com/en-us/azure/aks/use-multiple-node-pools), [GCP](https://cloud.google.com/kubernetes-engine/docs/concepts/node-pools) and [AWS](https://docs.aws.amazon.com/eks/latest/userguide/worker.html)). Each pool is a group of VMs with a specific spec/type. This allows a cluster to support, for example, both high memory VMs and GPU VMs. Workloads can then be scheduled onto these different VM types using [`nodeSelectors`](https://cloud.google.com/kubernetes-engine/docs/concepts/node-pools#deploying_services_to_specific_node_pools).  

A `nodeSelector` provides controls around which nodes each pod gets scheduled on. Typically developers don't include `nodeSelectors` which means that an admin needs to write a mutating admission controller to tell k8s to automatically add them and assign the workload to the correct type of machine in the cluster.

In this example, we show how to use an OPA Policy to assign `nodeSelectors` to ensure that `pods` from different k8s `namespaces` run on different `agentpools`.

The policy looks at the `namespace` in which the `pod` is to be scheduled and checks for an `agentpool` label. If the `namespace` has a label set then a `nodeselector` is set to assign the pod to the `agentpool`. 

For example, with the policy deployed the following `namespace`

```yaml
apiVersion: v1
kind: Namespace
metadata:
  labels:
    agentpool: pool1
  name: mapped
```

A `pod` being created in the `namespace: mapped` will have a `nodeselector` of `agentpool: pool1` injected:


```yaml
apiVersion: v1
  kind: Pod
  metadata:
    name: testpod212
    namespace: mapped
  spec:
    containers:
    - args:
      - /bin/sh
      - -c
      - while true;do date;sleep 5; done
      image: busybox
      imagePullPolicy: Always
      name: sleep
    nodeSelector:
      agentpool: pool1
```

Additionally the policy will enforce that, if a `pod` with a `nodeselector` already set is submitted to a `namespace` with an `agentpool` label it is rejected to ensure that users can't override the assignment. 

## Deployment

Use `make docker-e2e` to test locally. This will: 

1. Build a docker image with all the pre-reqs (OPA, KIND, OpenSSL and python)
2. Create a KIND Kubernetes cluster locally 
3. Deploy Open Policy Agent Gateway into the cluster
4. Deploy `main.rego` into the cluster
5. Use [the Pyton integration tests](./integration/test/) to validate the labels are applied correctly. 

To develop I recommend using the [VSCode Remote - Containers](https://code.visualstudio.com/docs/remote/containers) addin which will start a VSCode instance within the [devcontainer](./.devcontainer). This means all dependencies, addins and settings are ready to go. Once the addin is installed open the folder and `CTRL+SHIFT+P` -> `Reopen in container`. 

Once VSCode has built the container you can use: 

- `make kind-start` to create a test cluster
- `make kind-deploy` to deploy OPA and `main.rego` into the cluster
- `make integration` to run the integration tests against the cluster
- `make update-rego` to push an update to `main.rego`. `make force-update-rego` can be used to also restart the OPA server in case of an error. 
- `make logs` to get the logs from OPA in the cluster

## Folder structure

- `./devcontainer` contains the `dockerfile` to support [VSCode Remote - Containers](https://code.visualstudio.com/docs/remote/containers) and the `docker` based integration tests. 
- `./integration` contains the deployment scripts and integration tests
- `./testdata` contains useful samples for testing/developing the policy
- `authz.rego` is the policy to secure access to the OPA server
- `data_input.rego` is a mock set of data used for testing the policy
- `main.rego` is the policy applying the `nodeselector`
- `test_node_selector.rego` is a set of unit tests to ensure the policy works as expected
- `requirement.txt` tracks the python dependencies for the integration tests

## Security 

By default the OPA server runs with no authentication/authorization setup. This would allow other users of the cluster to talk to the REST API and make
changes to the policies. 

To avoid this the sample [deploys `authz.rego`](./authz.rego) which limits the anonymous access of the OPA API to only the decision endpoint (`POST /`) and the Health endpoint (`GET /HEALTH`).

When deployed in the cluster the `kube-mgmt` sidecar reads in the policies defined in `configMaps` and uses the API on the OPA Server to load the policies. This process has to be secured via a `token` to ensure only the `kube-mgmt` can edit polices. 

A `token` is generated in `/integration/deploy/deploy.sh` and injected into `authz.rego`. Both the token and the `authz.rego` output are stored in a secret in the cluster. This is mounted as a volume by the OPA server and the `kube-mgmt` sidecar. All request received by the OPA Server check the `token` (`- "/policies/authz.rego"`) and all requests sent by the `kube-mgmt` sidecar append the token (`"--opa-auth-token-file=/policies/token"`).

For more details see:
- https://github.com/open-policy-agent/kube-mgmt/blob/master/docs/admission-control-secure.md
- https://www.openpolicyagent.org/docs/latest/security/

## Update Policy

The `kube-mgmt` sidecar watches changes to the `configmaps` and loads new or updates existing policies during normal execution. If update issues are occurring review the sidecars logs with `make logs-sidecar`

## Want more mutating?

This sample was written to be simple. For more complex scenarios review the [OPA Library](https://github.com/open-policy-agent/library/tree/master/kubernetes/mutating-admission), it has a section with helpers for more advanced work. 