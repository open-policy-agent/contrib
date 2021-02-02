# Kubernetes authorization webhook using OPA

Runnable Kubernetes
[authorization webhook](https://kubernetes.io/docs/reference/access-authn-authz/webhook/)
example using OPA for authorization policy decisions.

## Running

### Prerequisites

* [kind](https://kind.sigs.k8s.io/)
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

### Setup

With those installed, simply run:

```shell
./setup.sh
```

The demo uses [Kind](https://kind.sigs.k8s.io/) to launch a local Kubernetes
cluster and then deploys OPA to that, with Kubernetes authorization policies
deployed from the [policy](policy) directory. Kind uses kubeadm so that's the
config format used for providing the authorization webhook flags to the API
server (see [kind-conf.yaml](#kind-conf.yaml).

### Testing the auhtorization webhook

With the cluster up and running, you may now issue the usual `kubectl` commands
to interact with your local Kubernetes API. Since the default user for kind is a
cluster admin with all priveleges granted it won't autmoatically be evaluated by
the authorizer webhook (as the RBAC module is configured in front of it). In
order to work around this, you could either setup a service account - or perhaps
easier; just simulate requests from other users by using the
[impersonation](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#user-impersonation)
feature of kubectl:

```shell
$ kubectl get pods \
        --namespace kube-system \
        --as=someuser \
        --as-group=system:authenticated \
        --as-group=devops

Error from server (Forbidden): OPA: denied access to namespace kube-system
```

The OPA server is configured to print decisions to stdout, so simply view the logs
of the OPA pod (in the `opa` namespace) to see requests and responses.

## Updating policy

Change the policy under the policy directory and run `kubectl apply -k .` Note
that it may take a while before the policy change is reflected in the running
system.

## Tests

There's a couple of end-to-end tests using kubectl to test authorization policy
enforcement in the `test.sh` script. Simply run it to have them executed:

```shell
$ ./test.sh
All tests successful
```

## Cleanup

```shell
kind delete cluster --name opa-authorizer
```
