# Image Enforcer

This directory shows how you can use OPA to enforce admission policies on
container images in Kubernetes.

You can use this project to enforce simple admission policies on images such as:

- Naming conventions
- Version pinning (e.g., no latest)
- Registry usage (e.g., production must use internal registry)

Using the `openpolicyagent/clair-layer-sync` image, you can also enforce more
sophisticated policies based on image vulnerability scanning.

This directory contains the following:

- Kubernetes manifests for deploying the image enforcer.
- Steps to use the image enforcer.
- Source for the `openpolicyagent/clair-layer-sync` Docker image.

## Testing

This repository contains sample policies and a data set that includes layer metadata with known vulnerabilities.

### Prerequisites

- OPA v0.4.10 or later

This example has been tested using:

- Kubernetes 1.6.0
- CoreOS Clair (quay.io/coreos/clair:latest)[1]

Some familiarity with bootstrapping Kubernetes is useful for this example as it relies on configuring the API server to enable the [ImagePolicyWebhook](https://kubernetes.io/docs/admin/admission-controllers/#imagepolicywebhook) admission controller. If you are using minikube, you will need to use (i) the `--extra-config` option and (ii) `mount` sub-command to set and load the necessary configuration.

[1] The CoreOS Clair installation steps link to the `latest` tag. The same
steps do not work with the most recent pinned version (`v2.0.0`).

### Steps

First, run OPA in interactive mode and provide the sample policies and data as input:

```shell
opa run ./policies ./data/* -w
```

- `./policies` contains Rego policies that authorize images
- `./data` contains Docker image metadata and CoreOS Clair vulnerability metadata
- `-w` tells OPA to reload any of the files if they change

The `./policies` directory contains mock input data. Dump the mock data:

```
> data.mock
```

You should see two keys: ``allowed`` and ``denied``. These keys represent the webhook calls that OPA would receive from the ImagePolicyWebhook admission controller in Kubernetes.

Set the input document to the ``allowed`` data set. Expressions that reference the `input` document will use the value defined by ``allowed``.

```
> package repl
> input = data.mock.allowed
```

Now, evaluate the image enforcer policy. The result should indicate the input is allowed:

```
> data.io.k8s.image_policy.verify
```

Next, switch the input document to the ``denied`` data set.

```
> unset input
> input = data.mock.denied
```

Finally, evaluate the image enforcer policy again. Tthe result should indicate the input is denied:

```
> data.io.k8s.image_policy.verify
```

That's it. For more detail, take a look at the policy files under `./policies`:

- `default.rego` implements the entry point to the policy (`data.io.k8s.image_policy.verify`).
- `scan.rego` implements the layer/vulnerability scanning policy.
- `mock/input.rego` defines the mock input used above.

## Deployment

### 1. Deploy the image enforcer on top of Kubernetes along with CoreOS Clair

Deploy the CoreOS Clair (from the docs):

```shell
git clone https://github.com/coreos/clair
cd clair/contrib/k8s
kubectl create secret generic clairsecret --from-file=./config.yaml
kubectl create -f clair-kubernetes.yaml
```

The initial CoreOS Clair deployment can take upwards of 30-60 minutes to
complete. It has to download vulnerability data from different sources and
store it in Postgres. Use `kubectl logs ${clair_pod_name}` to see when it's
done. If the deployment was successful, you'll eventually see something like
this:

```
2017-05-16 23:04:28.171978 I | updater: updating vulnerabilities
2017-05-16 23:04:28.172007 I | updater: fetching vulnerability updates
2017-05-16 23:04:28.172098 I | updater/fetchers/ubuntu: fetching Ubuntu vulnerabilities
2017-05-16 23:04:28.172284 I | updater/fetchers/debian: fetching Debian vulnerabilities
2017-05-16 23:04:28.172415 I | updater/fetchers/rhel: fetching Red Hat vulnerabilities
2017-05-16 23:12:22.131390 I | updater: adding metadata to vulnerabilities
2017-05-16 23:54:53.174039 W | updater: fetcher note: Ubuntu precise/esm is not mapped to any version number (eg. trusty->14.04). Please update me.
2017-05-16 23:54:53.174130 W | updater: fetcher note: Ubuntu zesty is not mapped to any version number (eg. trusty->14.04). Please update me.
2017-05-16 23:54:53.177667 I | updater: update finished
```

Grab a ☕ and wait for this to complete or skip ahead but keep in mind the vulnerability data
will not be ready yet.

Create a ConfigMap containing the initial policy to load into OPA.

```shell
kubectl create configmap image-enforcer --from-file ./policies/default.rego
```

Deploy the image enforcer. You can use the manifest from this directory to start.

```shell
kubectl create -f ./manifests/deployment.yaml
```

Lastly, create a service to expose OPA's API.

```shell
kubectl create -f ./manifests/service.yaml
```

If you query OPA or follow the logs, you'll see layer metadata being pushed in.

### 2. Re-configure the Kubernetes API server to enable the ImagePolicyWebhook admission controller

Create a kubeconfig file (`/srv/kubernetes/image_policy/kubeconfig`) to contact OPA:

```yaml
clusters:
- name: opa
  cluster:
    server: http://image-enforcer.default.svc.cluster.local:8181/v0/data/io/k8s/image_policy/verify
users:
- name: apiserver
  user:
    token: apiserver
current-context: webhook
contexts:
- name: webhook
  context:
    cluster: opa
    user: apiserver
```

Create the configuration file for the admission controller (`/srv/kubernetes/admission/config.json`):

```json
{
  "imagePolicy": {
    "kubeConfigFile": "/srv/kubernetes/image_policy/kubeconfig",
    "defaultAllow": false
  }
}
```

Re-configure the Kubernetes API server to start with:

```
--admission-control=ImagePolicyWebhook
--admission-control-config-file=/srv/kubernetes/admission/config.json
```

Re-start the Kubernetes API server.

## openpolicyagent/clair-layer-sync

This Docker image replicates data between:

- Image registries and CoreOS Clair (layer data)
- Image registries and OPA (layer metadata)
- CoreOS Clair and OPA (layer vulnerability metadata)

See [deployment.yaml](./manifests/deployment.yaml) for an example of how to deploy the container alongside OPA.
