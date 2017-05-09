# HTTP API Authorization with OPA

This directory helps provide fine-grained, policy-based control over who
can run which RESTful APIs.  A tutorial is available at
[openpolicyagent.org/examples/http-api-authorization/](http://www.openpolicyagent.org/examples/http-api-authorization/).


## Contents

* A sample web application that asks OPA for authorization before executing an API call (`docker/`)
* A default policy that allows `/finance/salary/<user>` for `<user>` and for `<user>`'s manager (`docker/policy`)


## Setup

The web application and OPA both run in docker-containers.  For convenience we
included a docker-compose file, so you'll want
[docker-compose](https://docs.docker.com/compose/install/) installed.

To build the containers and get them started, use the following make commands.

```
make       # build the containers with docker
make up    # start the containers with docker-compose
```
