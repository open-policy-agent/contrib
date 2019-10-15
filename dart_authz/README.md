# ![logo](logo.png) OPA API Authorization in Dart Example

This repository shows how to integrate a service written in Dart with [OPA](http://www.openpolicyagent.org) to perform
API authorization. It is a direct port of the [OPA-Python](https://github.com/open-policy-agent/example-api-authz-python)
example, with a few enhancements.

## Trying the example

This example utilizes an independent OPA server which must already be running, and which must allow new policies to be
uploaded. An existing OPA server URI can be defined in the `OPA_URL` environment variable, which otherwise defaults to
a local instance at `http://localhost:8181`.

To run the OPA instance locally:

```bash
$ opa run -s
```

note that the example policy (in the `./policies` directory) will be uploaded to the OPA server by the application
directly. Any additional policies with the `.rego` extension found in this directory will similarly be uploaded to the
OPA server on start.

Run the server:

```bash
$ dart bin/server.dart
```

Without authorization, view a list of cars:

```bash
$ curl -X GET localhost:8080/cars
```

As someone with the manager role, create a car (this should be allowed):

```bash
$ curl -H 'Authorization: alice' -H 'Content-Type: application/json' \
    -X PUT localhost:8080/cars/test-car \
    -d '{"model": "Toyota", "vehicle_id": "357192", "owner_id": "4821", "id": "test-car"}'
```

As someone with the car admin role, try to delete a car (this should be denied):

```bash
$ curl -H 'Authorization: kelly' \
    -X DELETE localhost:8080/cars/test-car
```

## Running from Docker

To run from Docker, simply specify the host and port of the OPA server through
the passed in `OPA_URL` environment variable:

```bash
$ docker run -e OPA_URL='opa:8181' -p 8080:8080 openpolicyagent/demo-dart:latest
```

Note that by default the Docker image enables the Dart Observatory, which binds
port 8181 within the container by default. If using `--net=host`, the default
Observatory port needs to be shifted out of the way. This can be done by
tweaking the `DART_VM_OPTIONS`, as so:

```bash
$ docker run -e DART_VM_OPTIONS='--enable-vm-service=8282' --net=host openpolicyagent/demo-dart:latest

Starting Dart with additional options --enable-vm-service=8282
Observatory listening on http://127.0.0.1:8282/4y7welzb8Fc=/
Applying policy: ./policies/example.rego
Example Service listening on 0.0.0.0:8080
...
```

## Features and bugs

Please file feature requests and bugs at the [issue tracker][tracker].

[tracker]: https://github.com/adaptant-labs/opa-api-authz-dart/issues

## License

Licensed under the terms of the Apache 2.0 license (the license under which the `OPA-Python` example was released),
the full version of which can be found in the [LICENSE](LICENSE)
file included in this distribution.

## Acknowledgements

- Derived from [example-api-authz-python](https://github.com/open-policy-agent/example-api-authz-python) by [@tsandall](https://github.com/tsandall).
