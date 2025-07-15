# AuthZEN Interop

`./policy` contains the rego policies and data to be loaded into OPA to pass the AuthZEN interop scenario tests.

```
opa run -s -b policy --addr http://localhost:8181
```
