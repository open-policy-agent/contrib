# kafka-authorizer

This directory contains a Kafka authorizer plugin that uses OPA to enforce
policy decisions.

**DEPRECATED** This authorizer plugin has been replaced by the [opa-kafka-plugin](https://github.com/anderseknert/opa-kafka-plugin),
which works on modern versions of Kafka, and provides additional features. The code here is kept here for historical reasons only.

For more information on integrating Open Policy Agent with Apache Kafka, see the
[OPA documentation](https://www.openpolicyagent.org/docs/latest/kafka-authorization/) on the topic.

## Build

To build the plugin, you must have Docker installed. If you want to build
without Docker `mvn install` should be enough but requires that you have Java
and Maven installed on your system.

Run `make build` to build and package the plugin.

The JARs required to run the plugin will written to `./target/kafka-authorizer-opa-<VERSION>-package/share/java/kafka-authorizer-opa`.

## Install

Include the following JAR files in the Kafka classpath:

* `kafka-authorizer-opa-<VERSION>.jar`
* `gson-2.8.2.jar`

Enable the plugin by adding the following to `server.properties`:

```
authorizer.class.name: com.lbg.kafka.opa.OpaAuthorizer
```

The plugin supports the following properties:

| Property Key | Example | Description |
| --- | --- | --- |
| `opa.authorizer.url` | `http://opa:8181/v1/data/kafka/authz/allow` | Name of the OPA policy to query. |
| `opa.authorizer.allow.on.error` | `false` | Fail-closed or fail-open if OPA call fails. |
| `opa.authorizer.cache.initial.capacity` | `100` | Initial decision cache size. |
| `opa.authorizer.cache.maximum.size` | `100` | Max decision cache size. |
| `opa.authorizer.cache.expire.after.ms` | `600000` | Decision cache expiry in milliseconds. |
| `opa.authorizer.token` | `` | Token for authentication with OPA. |

## Usage

The OPA policy should return a boolean value indicating whether the action
should be allowed (`true`) or denied (`false`).

The plugin provides input data describing the principal, operation, and
resource.

```ruby
# Example principle information.
input.session.principal.principalType = "User"
input.session.principal.name = "ANONYMOUS"
input.session.clientAddress = "127.0.0.1"
input.session.sanitizedUser = "ANONYMOUS"

# Example operation information.
input.operation.name = "ClusterAction"

# Example resource information.
input.resourceType.name = "Cluster"
input.resource.name = "kafka-cluster"
```

The following table summarizes the supported resource types and operation names.

| `input.resourceType.name` | `input.operation.name` |
| --- | --- |
| `Cluster` | `ClusterAction` |
| `Cluster` | `Create` |
| `Cluster` | `Describe` |
| `Group` | `Read` |
| `Group` | `Describe` |
| `Topic` | `Alter` |
| `Topic` | `Delete` |
| `Topic` | `Describe` |
| `Topic` | `Read` |
| `Topic` | `Write` |
