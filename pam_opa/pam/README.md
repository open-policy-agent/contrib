# Styra's PAM authorization module

This implements a PAM module that integrates with Open Policy Agent (OPA)
to make authorization decisions. When invoked, the module makes an authorizaton
decision with the help of OPA in three cycles:

1. *display*: The module asks OPA what should be displayed or prompted to the user,
   and collects the user responses.

2. *pull*: The module asks OPA what data (files and environment variables) should be pulled
   from the system, and collects the requested data.

3. *authz*: The module sends all collected data to OPA, and makes the authorization decision
   based on the response received.

## Building the module

To compile the authz PAM module, the following are needed:

1. [PAM development library](http://www.linuxfromscratch.org/blfs/view/svn/postlfs/linux-pam.html).
   On debian, this can be obtained using `apt-get install libpam0g-dev`.

2. [cURL library](https://curl.haxx.se/download.html). On debian, this can be obtained
   using `apt-get install libcurl4-gnutls-dev`

3. [jansson library](http://www.digip.org/jansson/).

Preferably, all libraries should be installed under `/usr`, so that the `LD_LIBRARY_PATH` environment
variable is not required at runtime. This can be done, for example, by running the autoconf
configuration files with the `prefix` flag during installation.

```shell
$ ./configure --prefix=/usr
```

Finally, build the module using `make`.

It is also possible to build the module in a docker container, which only requires that docker
be installed. See [this README for instructions](../README.md).

## Installation

1. Copy the module to `/lib/security` or wherever PAM modules reside on your
system.

2. Add the PAM module to an application's [PAM configuration](https://linux.die.net/man/5/pam.conf),
   for example add the following to `/etc/pam.d/sudo`:

```
auth required /lib/security/pam_opa.so url=http://opa:8181 authz_endpoint=/v1/data/sshd/authz display_endpoint=/v1/data/display pull_endpoint=/v1/data/pull log_level=debug
```

**Warning:** Once this PAM module is configured, OPA must be responsive and configured with the appropriate policies or commands configured to use the OPA PAM module will fail. The example above blocks using `sudo` unless the endpoint `http://opa:8181/v1/data/sshd/authz` returns the expected response, which would prevent removing the module with a non-root account since `/etc/pam.d/sudo` requires admin privileges to edit. In such a situation, load a policy that sets `allow := true` into OPA using its [policies management APIs](https://www.openpolicyagent.org/docs/latest/rest-api/#create-or-update-a-policy).

## Configuration

This section breaks down the different pieces in the PAM config `/etc/pam.d/sudo` example from above.

#### Type

`auth` is the PAM type. Only `auth` and `account` types are implemented in this module,
 and they perform exactly the same operations.

#### Control

`required` is the PAM control level, indicating that a failure code from this PAM
module should ultimately lead to a the application being sent a failure response.

#### Module and arguments

`lib/security/pam_opa.so` is the full path to the PAM module.
The module accepts arguments of the format `<flag>=<value>`.

###### Valid flags

|Property           |Required   |Description|
|-------------------|-----------|-----------|
|`url`              |yes        |The URL of an OPA instance API.
|`sock`             |no         |The path to a unix socket OPA is listening on.
|`display_endpoint` |no         |The path of the package containing policy that describes what to display or prompt the user.
|`pull_endpoint`    |no         |The path of the package containing policy that describes the JSON files and environment variables that should be collected from the system.
|`authz_endpoint`   |yes        |The path of the package containing the policy that takes all collected data as input and makes the final decision.
|`log_level`        |no         |The verbosity of logs that this PAM module generates.

While not providing non-required endpoints will not break the PAM module, recommended practice is
to have valid endpoints in all PAM configurations.

This will ensure that configurations only need to be modified once, requiring minimal provisioning in the future.
Policy can then control all authorization behavior. For example, to remove all user interaction
from the process, simply have the *display* policy (described below) evaluate to an empty list.

The same *display*, *pull* and *authz* packages have been used for both the *sudo* and *sshd* in this document.
In production, it is more useful to use separate, fine-grained policy packages in each PAM configuration file.

## Policies

Requirements and examples of the OPA packages for each cycle are described below.
Note that it is OK to do nothing in the *display* and *pull* cycles: the associated policies should then evaluate to empty lists.

#### Display

The only rule required for this cycle is `display_spec`, which must contain a list of objects.
These objects each describe a message that should be displayed to the user, in order.
Each object should contain:

1. `message`: The message to display to the user.
2. `key`: The key that the user's response should be associated with, where applicable.
    The authorization policy will ultimately be invoked with an object containing this key.
3.  `style`: One of the following PAM conversation styles:
    - `prompt_echo_on` prompts the user for non-sensitive information.
    - `prompt_echo_off` prompts the user for sensitive information.
    - `info` displays an informational message to the user.
    - `error` displays an error message to the user.

    The actual conversation between the application and the user is implemented by the
    application, and may vary in behavior. For example, some versions of OpenSSH will postpone
    displaying all `info` messages, dumping them at the end after all prompts are completed.

    Each application has a different maximum input length that the user can enter.
    This value is, for example, 256 characters for common implementations of `sudo`, and
    1024 characters for OpenSSH.
    Input larger than this maximum will usually be truncated by the application.

###### Example

The following policy greets the user, and then prompts the user for their last name and secret.

```
# This package path should be passed with the display_endpoint flag
# in the PAM configuration file.
package display

display_spec = [
	{
		"message": "Welcome to the OPA-PAM demonstration.",
		"style": "info",
	},
	{
		"message": "Please enter your last name: ",
		"style": "prompt_echo_on",
		"key": "last_name",
	},
	{
		"message": "Please enter your secret: ",
		"style": "prompt_echo_off",
		"key": "secret",
	},
]
```


#### Pull

The following rules are required:

1. `files` should be a list of strings, each being a path to a JSON file on the system
   that the PAM module is running on. Only absolute paths are guaranteed to work.

2. `env_vars` should be a list of strings, each being a name of an environment variable
   whose value is needed for authorization. The environment variable should be readable
   by the PAM module.

###### Example

Let's assume that we have several running hosts each having a file `/etc/host_identity.json`
which looks like this -
```
{
    "host_id": "<some host id>"
}
```
The following policy requests for collection of the JSON file's contents.

```
# This package path should be passed with the pull_endpoint flag
# in the PAM configuration file.
package pull

# JSON files to pull.
files = ["/etc/host_identity.json"]

# env_vars to pull.
env_vars = []
```

#### Authz

The following rules are required:

1. `allow` should evaluate to `true` if the authorization is successful.
2. `errors` should be an array containing error messages that the PAM module will log.

The authz package will receive an `input` object containing the data for making the decision:
- `display_responses` is an object having keys as defined in the display policy, and
  user responses to prompts as values.
- `pull_responses` is an object containing the following:
  - `files` is an object having file paths as keys and their contents as values.
  - `env_vars` is an object containing environment variable names and values.
- `sysinfo` is an object containing the following default system information, extracted
  from the PAM session.
  - `pam_username` is the username that the session will grant after authorization.
  - `pam_service` is the name of the application which invoked the PAM session.
  - `pam_req_username` is the username that made the authorization request.
  - `pam_req_hostname` is the hostname that made the authorizatoin request.

###### Example

The two previous cycles determine what the *authz* policy receives as context.
Based on the *display* and *pull* examples above, the context should be:

```
{
    "input": {
        "display_responses": {
            "last_name": "<user input>",
            "secret":    "<user input>"
        },
        "pull_responses": {
            "files": {
                "/etc/host_identity.json": {
                    "host_id": "<some host id>"
                }
            },
            "env_vars": {}
        },
        "sysinfo": {
            "pam_username":     "<PAM session value>",
            "pam_service":      "<PAM session value>",
            "pam_req_username": "<PAM session value>",
            "pam_req_hostname": "<PAM session value>"
        }
    }
}
```

The example policy below only grants access if

- The user enters `ramesh` and `suresh` when prompted.
- The file `/etc/host_identity.json` has `"host_id": "frontend"`.
- The username requesting authorization is `ops`.

```
# This package path should be passed with the authz_endpoint flag
# in the PAM configuration file.
package sshd.authz

import input.display_responses
import input.pull_responses
import input.sysinfo

default allow = false

allow {
	# Verify user input.
	display_responses.last_name = "ramesh"
	display_responses.secret = "suresh"

	# Only allow running on host_id "frontend"
	pull_responses.files["/etc/host_identity.json"].host_id = "frontend"

	# Only authorize user "ops"
	sysinfo.pam_username = "ops"
}

errors["You cannot pass!"] {
	not allow
}
```

## Logging

The PAM module logs using syslog to the `LOG_AUTH` facility. On Linux, the logs can
usually be found at `/var/log/auth.log`.

The following log levels are accepted with the `log_level` flag:

- `none`: Don't log anything.
- `error`: Log only error messages.
- `info`: Log error and info messages.
- `debug`: Log verbosely; additionally log to standard error.

## Debugging

Set up the `sudo` and `sshd` PAM configuration files to run with `log_level=debug`.
It is recommended to use `sudo` for debugging first, instead of `sshd`.

*Note:* The
[trial docker containers](../README.md#try-it),
by default, run with settings suitable for debugging.

To debug a docker container running these PAM configurations, get into the container,

```shell
$ docker exec -it <container-id> bash
```

Then switch to the user you want to test with, for example `ops`.

```shell
$ su - ops
```

Now run `sudo ls` and read the debug logs.

If you want to additionally debug `sshd`, it is recommended to start your own
instance of `sshd` in debug mode because the existing `sshd` daemon will not log to stderr.

```shell
$ $(which sshd) -d -p 2227
```

To SSH as user `ops` to your `sshd` service, run the SSH client with lenient
host requirements and verbose logs.

```bash
ssh -p 2227 ops@<container-ip> -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -vvv
```

## Security

OPA should be secured so that non privileged users can not contact OPA and change policy or use it to help break in. This module does not yet support any authentication but when used in conjunction with the -sock flag, security can be enforced with unix file permissions (chown root.root opa.sock and chmod 600 opa.sock). Start the OPA server with -a unix://<path>/opa.sock.
