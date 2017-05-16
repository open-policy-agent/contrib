# Open Policy Agent - Contributions

This repository holds integrations, examples, and proof-of-concepts that work with the Open Policy Agent (OPA) project.

## Examples and Integrations

- [HTTP/micro-service API authorization with Linkerd](./linkerd_authz)
- [HTTP API Authorization Example (Python)](./api_authz)
- [PAM module - SSH and sudo Authorization](./pam_authz)
- [Puppet Authorization Example](./puppet_example)
- [Kubernetes Image Policy Enforcer](./image_enforcer)

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
