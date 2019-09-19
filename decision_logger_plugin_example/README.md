## Custom Decision Logger Plugin

This directory contains an example of implementation a custom decision
logger for OPA.

## Build

```bash
go build
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
./decision_logger_plugin_example run --server --config-file=config.yaml
```

Exercise the custom decision logger:

```bash
curl localhost:8181/v1/data
```

Example output. The 4th line is from the println logger. The other lines are standard OPA log messages.

```
{"addrs":[":8181"],"insecure_addr":"","level":"info","msg":"Initializing server.","time":"2019-09-19T11:49:16-04:00"}
{"level":"info","msg":"Starting decision logger.","plugin":"decision_logs","time":"2019-09-19T11:49:16-04:00"}
{"client_addr":"127.0.0.1:41150","level":"info","msg":"Received request.","req_id":1,"req_method":"GET","req_path":"/v1/data","time":"2019-09-19T11:49:21-04:00"}
{map[id:280bcbde-9195-4142-a055-d9661d4805fc version:] ad492980-c1d2-415a-bb6a-71cba6111e59  map[] data  <nil> 0xc000331a10 [] <nil> 127.0.0.1:41150 2019-09-19 15:49:21.561760024 +0000 UTC map[timer_rego_input_parse_ns:1273 timer_rego_load_bundles_ns:635 timer_rego_load_files_ns:2173 timer_rego_module_parse_ns:659 timer_rego_query_compile_ns:123733 timer_rego_query_eval_ns:24995 timer_rego_query_parse_ns:724709 timer_server_handler_ns:931606]}
{"client_addr":"127.0.0.1:41150","level":"info","msg":"Sent response.","req_id":1,"req_method":"GET","req_path":"/v1/data","resp_bytes":66,"resp_duration":1.780271,"resp_status":200,"time":"2019-09-19T11:49:21-04:00"}
```
