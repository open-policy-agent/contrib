# Puppet Authorization Example

This example shows how OPA can authorize changes to Puppet manifests.

You can try out the steps below using the policy file and JSON data contained
in this repository.

This example works by defining a policy that:

1. Accepts a compiled Puppet manifest as input (the Puppet catalog)
1. Scans the resources in the Puppet catalog and identifies the author (using Git blame data)
1. Checks if the author is a member of the team responsible for the resource
1. Produces a `true` or `false` result that indicates if the Puppet catalog is allowed (based on the checks above)

The policy relies on two pieces of external data:

1. Git blame identifying the Puppet resource author
1. Team membership

In real-world scenarios, this data could be replicated into OPA from sources like GitHub and LDAP.

## Steps

### 1. Start OPA and load the example policy and data.

First, start OPA using Docker.

```bash
docker run -it --rm -p 8181:8181 openpolicyagent/opa:0.4.8 run -s -l debug
```

Next, load the example policy:

```bash
curl -X PUT --data-binary @puppet_authz.rego localhost:8181/v1/policies/puppet_authz
```

Finally, load an example Git blame data set.

```bash
curl -X PUT -d @blame_allowed.json localhost:8181/v1/data/git
```

### 2. Test the policy with a query.

Now, run a policy query against OPA, providing the example Puppet catalog as input.

```bash
curl -X POST -d @puppet_catalog.json localhost:8181/v1/data/puppet/authz/allow
```

The result indicates that `puppet_catalog.json` is allowed:

```http
HTTP/1.1 200 OK
Content-Type: application/json
```

```json
{
    "result": true
}
```

### 3. Update the Git blame data in OPA.

Let's update the Git blame data in OPA.

```bash
curl -X PUT -d @blame_denied.json localhost:8181/v1/data/git
```

In this Git blame data set, `bob` (who is a member of `app_team`) modified an `infra_team` resource.

### 4. Test the policy with another query.

```bash
curl -X POST -d @puppet_catalog.json localhost:8181/v1/data/puppet/authz/allow
```

This is the same query that we ran before. However, this time, the Git blame data stored inside OPA is different (and it represents a policy violation).  As a result, the catalog is denied.

```http
HTTP/1.1 200 OK
Content-Type: application/json
```

```json
{
    "result": false
}
```
