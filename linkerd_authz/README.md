# Micro-service Authorization with Linkerd

An experimental, policy-enabled, linkerd identifier that enforces authorization decisions.

## Building

```
./sbt assembly  # puts JAR under target/scala-2.11/
```

## Running (local)

Download and untar linkerd-0.9.0, change into directory.

1. `export L5D_HOME=$PWD` sets environment variable that linkerd-0.9.0-exec script uses to load plugins.
1. `mkdir -p plugins` creates directory to stick plugin JARs into.
1. Copy JAR file into `plugins` directory.
1. Create linkerd configuration at `config/opa_linkerd_example.yaml`:

    ```yaml
    namers:
    - kind: io.l5d.fs
    rootDir: disco

    routers:
    - protocol: http
    dtab: |
      /dns-hostport => /$/inet;
      /dns          => /$/io.buoyant.hostportPfx/dns-hostport | /$/inet;
      /pool         => /#/io.l5d.fs;
      /svc          => /pool | /dns;
    httpAccessLog: /dev/stdout
    identifier:
      kind: org.openpolicyagent.linkerd.authzIdentifier
      ip: 127.0.0.1
      port: 8181
      path: /v1/data/example/linkerd/authz
    servers:
      - port: 4140
        ip: 0.0.0.0
        addForwardedHeader:
          by:
            kind: "ip:port"
          for:
            kind: "ip"
    ```

1. Run linkerd:

    ```bash
    ./linkerd-0.9.0-exec config/opa_linkerd_example.yaml
    ```

## Manual Testing (local)

1. Start a simple webserver:

    ```bash
    python -m SimpleHTTPServer 9999
    ```

1. Start OPA:

    ```bash
    docker run -p 8181:8181 openpolicyagent/opa:0.4.6 run --server --log-level=debug
    ```

1. Define simple policy (`example.rego`):

    ```ruby
    package example.linkerd.authz

    import input.method
    import input.path

    errors["request denied by administrative policy"] {
      not allow
    }

    default allow = false

    allow {
      method = "GET"
      not contains(path, "deadbeef")
    }
    ```

1. Push policy into OPA:

    ```bash
    curl -X PUT localhost:8181/v1/policies/test --data-binary @example.rego
    ```

    and then, optionally, watch the policy for changes and push into OPA:

    ```
    fswatch -o example.rego | xargs -n1 \
      curl -X PUT localhost:8181/v1/policies/test --data-binary @example.rego
    ```

1. GET some document from webserver via linkerd:

    ```bash
    curl -H 'Host: web' localhost:4140
    ```

1. Try to POST some document to webserver via linkerd:

    ```bash
    curl -d 'hooray' -H 'Host: web' localhost:4140
    ```

    Response:

    ```
    Unknown destination: Request("POST /", from /127.0.0.1:59508) / request denied by administrative policy
    ```

That's it! 🎉

## More Information and Future Work

By default, the identifier is designed to **fail-closed**. This means that if
an error occurs while communicating with OPA or OPA returns an error, the
identifier will deny the request.

The identifier is currently hard-coded to use the
`io.buoyant.router.http.HeaderIdentifier` as the actual underlying identifier.
The `Host` header is used as input to the `HeaderIdentifier`.

Future work could include:

- Separating the authorization and identification steps so that underlying identifiers do not have to be wrapped.
- Adding support for sending the entire HTTP message body to OPA (in addition to headers, path, and method).
- Adding support for protocols other than HTTP.
