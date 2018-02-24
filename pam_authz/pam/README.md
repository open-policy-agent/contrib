# Styra's PAM authorization module

This implements a PAM module that integrates with Open Policy Agent (OPA)
to make authorization decisions. When invoked, the module queries OPA
and if no errors are returned, then the action is allowed. This allows
administrators to implement fine-grained access control.

This module is not intended to perform authentication. One can use standard
linux accounts for that or a user database such as LDAP.


## Usage

This repo includes a detailed [example usage of the pam_authz module](http://www.openpolicyagent.org/docs/ssh-and-sudo-authorization),
so the instructions here focus on the installation and configuration of the
PAM module. Details of writing policy to make the authorization decisions are
found in the docker example.

### Building the module

To compile the authz PAM module, the following are needed:

1. [golang compiler](https://golang.org/doc/install) (at least version 1.7)

2. PAM development library. On debian, this can be obtained using
   `apt-get install libpam0g-dev`.

Finally, build the module using `make`. This will automatically obtain the go dependencies.

It is also possible to build the module in a docker container, which only requires that docker
be installed. See [this README for instructions](../README.md).

### Installation

1. Copy the module to `/lib/security` or wherever PAM modules reside on your
system.

2. Configure it as an auth or account module, for example, add the following
line to `/etc/pam.d/sudo`:

```
account required /lib/security/pam_authz.so url=http://localhost:8181 policy_path=/v1/data/ssh/authz
```

Configuration options:
* `url` - the address of OPA. Defaults to `http://localhost:8181`.

* `policy_path` - the path to the policy endpoint used for authorization
  decisions. Required.

* `identity_file_path` - path to a json file. The contents of this file are sent to OPA
  to be used in the access control policy. For example, this file might contain a single
  string `"host-uuid"` that indicates the server's ID. Or it could contain an object with
  multiple fields. The intent is for this information to be used in the authorization
  policy to make allow or deny decisions. The default value is `/etc/host_identity.json`.

When pam_authz is invoked, it queries OPA the endpoint at the specified url.
The input document to this query contains the following fields: `user` and
`host_identity`. The `user` field is a string, and `host_identity` is required
to be json and is the contents of the `indentity_file_path` file described above.
The OPA policy needs to use this data to make an authorization decision.
Example policy documents are [included in this project](../docker/policy/ssh.rego).

The response from OPA is expected to include a field `errors`, which is an
array of strings. If `errors` is empty, then the module returns success.
Otherwise, the module returns deny.


## Copyright

Portions of this module use code from the [pam-ussh project](https://github.com/uber/pam-ussh).
Those portions are copyright Uber Technologies. See the [license](LICENSE)
for more details.
