# otlp_server_logger

An OPA contrib plugin that forwards OPA's internal operational logs to an OpenTelemetry collector via OTLP/gRPC or OTLP/HTTP.

> **Scope:** This plugin captures logs written through OPA's internal `m.logger` path — startup messages, bundle activity, auth decisions, HTTP request handling, and similar runtime events. It can also optionally receive decision log events (policy evaluation results) by setting `decision_logs.plugin: otlp_server_logger` in OPA config — OPA converts each `EventV1` to a structured `slog` record before forwarding. See [How it works](#how-it-works) for details.

Related issue: [open-policy-agent/opa#8214](https://github.com/open-policy-agent/opa/issues/8214)

## Prerequisites

- Go 1.25+
- OPA v1.15+ (`server.logger_plugin` support, introduced in OPA v1.15.0)
- An OpenTelemetry collector reachable from OPA (e.g. OTEL Collector, Grafana Alloy, Jaeger, etc.)

## Configuration

> **Activation:** OPA does not load this plugin automatically. You must set `server.logger_plugin: otlp_server_logger` in your OPA config to redirect all server logs through the plugin.

All fields live under `plugins.otlp_server_logger` in `opa.yaml` / `config.yaml`.

| Field | Type | Default | Description |
|---|---|---|---|
| `service` | string | `""` | Name of an OPA service (from the `services` section) to inherit connection settings from |
| `type` | string | inherited or **required** | Transport: `grpc` or `http` (inherited from `distributed_tracing` if set; otherwise required) |
| `address` | string | `localhost:4317` (gRPC) / `localhost:4318` (HTTP) | Collector endpoint (host:port, no scheme) |
| `service_name` | string | `opa` | Value of the OTel `service.name` resource attribute |
| `encryption` | string | `off` | TLS mode: `off`, `tls` (one-way), or `mtls` (mutual) |
| `tls_cert_file` | string | `""` | Client certificate PEM file (required for `mtls`) |
| `tls_private_key_file` | string | `""` | Client private key PEM file (required for `mtls`) |
| `tls_ca_cert_file` | string | system pool | CA certificate PEM file (optional; falls back to system trust store if omitted) |
| `level` | string | `info` | Minimum log level: `debug`, `info`, `warn`, or `error` |
| `headers` | map | `{}` | Extra HTTP/gRPC metadata headers (e.g. `Authorization: Bearer <token>`) |
| `compression` | string | `""` | Wire compression: `gzip` or `none` (empty means no compression) |
| `export_timeout_ms` | int | SDK default (~10s) | Per-export network timeout in milliseconds |

## Reusing existing OPA config

The plugin avoids duplicating connection settings by inheriting from OPA's existing config.
Priority order: explicit plugin fields → `service` reference → `distributed_tracing` fallback.

### Via `services` (recommended)

If the OTLP collector is already defined as an OPA service, reference it by name:

```yaml
services:
  otel-collector:
    url: https://otel-collector.monitoring.svc.cluster.local:4317
    credentials:
      bearer:
        token: my-token

server:
  logger_plugin: otlp_server_logger

plugins:
  otlp_server_logger:
    service: otel-collector  # inherits address, encryption, TLS, and auth headers
    type: grpc               # required: port numbers don't identify the protocol
    level: info
```

The plugin extracts `address`, `encryption`, TLS cert/key paths, and static auth headers
from the service config. `type` must always be set explicitly — port numbers do not
reliably identify the OTLP transport protocol. Dynamic credential sources (OAuth2, token
files) are not supported and must be set via `headers` instead.

### Via `distributed_tracing`

If OPA's `distributed_tracing` is already configured, connection fields are inherited
automatically when `service` is not set:

```yaml
distributed_tracing:
  type: grpc
  address: otel-collector.monitoring.svc.cluster.local:4317
  encryption: "off"

server:
  logger_plugin: otlp_server_logger

plugins:
  otlp_server_logger:
    level: info
```

Explicit values in `plugins.otlp_server_logger` always take precedence over inherited ones.

## Full OPA config example

Server logs only:

```yaml
plugins:
  otlp_server_logger:
    type: grpc
    address: otel-collector.monitoring.svc.cluster.local:4317
    service_name: opa
    encryption: "off"
    level: info

server:
  logger_plugin: otlp_server_logger
```

Server logs **and** decision logs via one plugin:

```yaml
plugins:
  otlp_server_logger:
    type: grpc
    address: otel-collector.monitoring.svc.cluster.local:4317
    service_name: opa
    encryption: "off"
    level: info

server:
  logger_plugin: otlp_server_logger

decision_logs:
  plugin: otlp_server_logger
```

With mTLS:

```yaml
plugins:
  otlp_server_logger:
    type: grpc
    address: otel-collector.monitoring.svc.cluster.local:4317
    service_name: opa
    encryption: mtls
    tls_cert_file: /etc/opa/tls/client.crt
    tls_private_key_file: /etc/opa/tls/client.key
    tls_ca_cert_file: /etc/opa/tls/ca.crt
    level: debug

server:
  logger_plugin: otlp_server_logger
```

With compression, timeout, and auth header:

```yaml
plugins:
  otlp_server_logger:
    type: http
    address: ingest.example.com:4318
    service_name: opa
    encryption: tls
    compression: gzip
    export_timeout_ms: 5000
    headers:
      Authorization: "Bearer <token>"
    level: info

server:
  logger_plugin: otlp_server_logger
```

> **Security note:** Values set under `headers` (including `Authorization`) are stored in
> OPA's resolved config and will appear in plaintext in the `/v1/config` API response.
> To avoid token exposure, prefer referencing a `service` that carries credentials — tokens
> inherited via `service:` are resolved at startup and are not included in the raw config JSON.

## Build

```sh
make build          # compiles and tests all packages
go build -o opa-otlp-server ./examples/  # produces the OPA binary
make test           # runs unit tests
make push           # builds & pushes multi-arch Docker image
```

To run locally (requires a reachable OTLP collector at the configured address):

```sh
go build -o opa-otlp-server ./examples/
./opa-otlp-server run --server --config-file examples/config.yaml

# Smoke test
curl localhost:8181/v1/data
```

## How it works

1. OPA calls `Factory.Validate` at startup to parse and validate the config.
2. OPA calls `Plugin.Start`, which creates an OTLP exporter and an OTel `LoggerProvider` with the following resource attributes automatically attached to every log record:

   | Attribute | Source |
   |---|---|
   | `service.name` | `service_name` config field (default: `opa`) |
   | `service.version` | OPA build version (e.g. `v1.4.0`) |
   | `host.name` | OS hostname, auto-detected |
   | `telemetry.sdk.name` / `.language` / `.version` | OTel Go SDK metadata, auto-detected |

3. The plugin wraps both a `slog.JSONHandler` (stdout) and `otelslog.NewHandler` (OTLP) in a `multiHandler`, then applies a `leveledHandler` that gates records by the configured minimum level. Every log record is written to stdout and exported to the OTLP endpoint simultaneously, so `kubectl logs` remains available as a fallback even when the collector is unreachable.
4. OPA calls `Plugin.Logger()` to retrieve the `slog.Handler` and uses it for all internal log output.
5. If `decision_logs.plugin: otlp_server_logger` is also set, OPA converts each decision event (`EventV1`) to structured `slog` attributes and forwards it through the same handler. Each decision log record carries fields such as `decision_id`, `path`, `input`, `result`, `timestamp`, and others.
6. Each log record is batched and exported to the configured OTLP endpoint.
7. On `Stop`, the provider is shut down gracefully (pending records are flushed before exit).
8. If OPA's built-in `services[].logs` OTLP export is also configured, both pipelines receive logs independently — this plugin does **not** replace that mechanism and running both will produce duplicate records at the collector.
9. `Reconfigure` inspects which fields changed:
   - **Hot-reload (no restart):** a `level`-only change atomically updates the minimum log level; the exporter stays running.
   - **Full restart (Stop + Start):** any change to `type`, `address`, `service_name`, `encryption`, TLS file paths, `compression`, `headers`, or `export_timeout_ms` triggers a graceful shutdown followed by a fresh `Start`. On restart failure, the previous config is restored and a recovery `Start` is attempted. If recovery also fails, the plugin status is set to `err`.
