# AuthZEN OPA Proxy

This repository contains a proxy server in `node.js` that exposes an AuthZen compliant API
while proxying requests and responses to an Open Policy Agent instance.

```
OPA_URL=<http://localhost:8181> OPA_POLICY_PATH=<authzen/allow> npm install && npm start
```
