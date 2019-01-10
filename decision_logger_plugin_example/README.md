## Custom Decision Logger Plugin

This directory contains an example of implementation a custom decision
logger for OPA.

## Build

```bash
go build -buildmode=plugin -o=plugin.so main.go
```

## Run

Create an OPA configuration file:

```yaml
decision_logs:
	plugin: println_decision_logger
plugins:
	println_decision_logger:
		stderr: false
```

Run OPA:

```bash
opa --plugin-dir=. run --server --config-file=config.yaml
```

Exercise the custom decision logger:

```bash
curl localhost:8181/v1/data
```

Example output:

```
INFO[2019-01-09T15:58:10-08:00] First line of log stream.                     addrs="[:8181]" insecure_addr=
INFO[2019-01-09T15:58:10-08:00] Starting decision log uploader.               plugin=decision_logs
INFO[2019-01-09T15:58:22-08:00] Received request.                             client_addr="127.0.0.1:51812" req_id=1 req_method=GET req_params="map[]" req_path=/v1/data
{map[id:d0dc7534-7b3f-42ee-8177-552d51a9504b] 09dd7698-f01a-4f75-891e-622db480ab0d  data 0xc0001da268 0xc0003c39b0 127.0.0.1:51812 2019-01-09 23:58:22.591477037 +0000 UTC}
INFO[2019-01-09T15:58:22-08:00] Sent response.                                client_addr="127.0.0.1:51812" req_id=1 req_method=GET req_path=/v1/data resp_bytes=66 resp_duration=1.63197 resp_status=200
```
