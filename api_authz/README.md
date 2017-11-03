# HTTP API Authorization with OPA

This directory helps provide fine-grained, policy-based control over who
can run which RESTful APIs.

A tutorial is available at [HTTP API Authorization](http://www.openpolicyagent.org/docs/http-api-authorization)


## Contents

* A sample web application that asks OPA for authorization before executing an API call (`docker/`)
* A default policy that allows `/finance/salary/<user>` for `<user>` and for `<user>`'s manager (`docker/policy`)
    * There are two policies given. The first is `api_authz.rego`, which is the default policy. The second is
      `api_authz_token.rego`, which allows you to perform the same task, but by communicating information relevant
      to the policy via JSON Web Tokens. The tokens to use for the second policy can be found in the `tokens`
      directory. Files with the `jwt` extension are the tokens themselves, and files with the `txt` extension
      are their respective decoded tokens for reference.

## Setup

The web application and OPA both run in docker-containers.  For convenience we
included a docker-compose file, so you'll want
[docker-compose](https://docs.docker.com/compose/install/) installed.

To build the containers and get them started, use the following make commands.

```
make       # build the containers with docker
make up    # start the containers with docker-compose
```

To instead use the example with JSON Web Tokens, use the following make commands.

```
make             # build the containers with docker
make up-token    # start the containers with docker-compose
```
