# Open Policy Agent - Contributions

This repository holds integrations, examples, and proof-of-concepts that work with the Open Policy Agent (OPA) project.

## Examples and Integrations

- [Kafka Authorization](./kafka_authorizer)
- [HTTP API Authorization (Spring Security)](./spring_authz)
- [HTTP API Authorization (Linkerd)](./linkerd_authz)
- [HTTP API Authorization (Python)](./api_authz)
- [HTTP API Authorization (Dart)](./dart_authz)
- [HTTP API Authorization (Kong)](./kong_api_authz)
- [SSH and sudo Authorization (PAM)](./pam_opa)
- [Puppet Authorization](./puppet_example)
- [Container Image Policy (Kubernetes and CoreOS Clair)](./image_enforcer)
- [Data Filtering (SQL)](./data_filter_example)
- [Data Filtering (Elasticsearch)](./data_filter_elasticsearch)
- [Data Filtering (MongoDB)](./data_filter_mongodb)
- [Data Filtering (Azure)](./data_filter_azure)
- [Cloud Foundry Policies](./cloud_foundry)
- [Decision Logger Plugin](./decision_logger_plugin_example)
- [IPTables (Linux)](./opa-iptables)
- [IPTables (Envoy)](./envoy_iptables)
- [JUnit Test Format Conversion](./junit)
- [Kubernetes Authorization](./k8s_authorization)
- [Kubernetes Node Selector](./k8s_node_selector)
- [Kubernetes API Client](./k8s_api_client)

This a list of integrations that are not currently part of the contrib repository:

- [HTTP API Authorization with Kong (built by @maxTN)](https://github.com/TravelNest/kong-authorization-opa)
- [Helm policies (built by @eicnix)](https://github.com/eicnix/helm-opa)
- [Install and load OPA with Habitat (built by @srenatus)](https://github.com/habitat-sh/core-plans/tree/master/opa)
- [Python/Flask extension for HTTP API Authorization and Policy Enforcement Points (built by @EliuX)](https://github.com/EliuX/flask-opa)
- [ANTLR v4 Rego Grammar (built by @anadon)](https://github.com/antlr/grammars-v4/tree/master/rego)
- [HTTP API Authorization with Slim Framework (PHP) (built by @segrax)](https://github.com/segrax/opa-php-examples)

## Contributing

If you have built an integration, example, or proof-of-concept on top of OPA that you would like to release to the community, feel free to submit a Pull Request against this repository. Please create a new top-level directory containing:

- A README.md explaining what your integration does
- A Makefile to build your integration

## Building and Releasing

Most integrations include a top-level Makefile with two targets:

* `build` - compiles/lints/tests the integration
* `push` - builds the integration and publishes artifacts

Many of the integrations produce one or more Docker images. These Docker images can be pushed to the hub.docker.com/u/openpolicyagent repository (assuming you are authorized.)

The Makefile in this directory contains `build` and `push` targets to build and push all integrations.
