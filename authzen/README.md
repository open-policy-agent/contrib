# AuthZEN Interop with OPA

This directory contains the components for running OPA as an [AuthZEN](https://openid.github.io/authzen/)-compliant authorization service.

The proxy translates between the AuthZEN evaluation API and OPA's data API so that any AuthZEN-compatible client can use OPA as a backend without modification.

## Contents

- [authzen-proxy](./authzen-proxy) - Node.js proxy server that exposes the AuthZEN evaluation endpoint and forwards requests to OPA
- [authzen-interop](./authzen-interop) - Rego policies and test data for the AuthZEN interop scenario tests

## Running

Start OPA with the interop policies:

```bash
cd authzen-interop
opa run -s -b policy --addr http://localhost:8181
```

In a separate terminal, start the proxy:

```bash
cd authzen-proxy
npm install
OPA_URL=http://localhost:8181 OPA_POLICY_PATH=authzen/allow npm start
```

The proxy listens on port `3000` by default. Send an AuthZEN evaluation request:

```bash
curl -X POST http://localhost:3000/access/v1/evaluation \
  -H "Content-Type: application/json" \
  -d '{
    "subject": { "id": "CiRmZDA2MTRkMy1jMzlhLTQ3ODEtYjdiZC04Yjk2ZjVhNTEwMGQSBWxvY2Fs" },
    "action": { "name": "can_read_todos" },
    "resource": { "type": "todo" }
  }'
```

## Configuration

The proxy accepts the following environment variables:

| Variable | Default | Description |
| --- | --- | --- |
| `PORT` | `3000` | Port for the proxy server |
| `OPA_URL` | `http://localhost:8181` | OPA server URL |
| `OPA_POLICY_PATH` | `authzen/allow` | OPA policy path to evaluate |

## Policy

The included Rego policy (`authzen-interop/policy/authzen/authzen.rego`) implements role-based access control with the following roles:

- **admin** - can create and delete resources
- **editor** - can create, update (own), and delete (own) resources
- **evil_genius** - can update any resource

Reading users and todos is allowed for all authenticated users regardless of role.
