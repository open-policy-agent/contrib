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
- [Grafana Dashboard](./grafana-dashboard)
- [OpenAPI Specification for OPA](./open_api)
- [SonarCloud Test Coverage Conversion](./sonarcloud)

For a comprehensive list of integrations, see the OPA [ecosystem](https://www.openpolicyagent.org/docs/latest/ecosystem/) page.

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
