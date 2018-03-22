# Policy-driven SSH and sudo with OPA and a PAM module

This directory helps provide fine-grained, policy-based control over who can
ssh and sudo into each of your servers and containers.

You can find a step-by-step tutorial at [SSH and sudo Authorization](http://www.openpolicyagent.org/docs/ssh-and-sudo-authorization).


### Directory contents

This directory includes:

* A policy-enabled PAM module that you install on each of your servers or containers ([/pam](./pam))
* Code showing you how to install and configure the PAM module and package servers as containers ([/docker](./docker))

### Try it

To get started, make sure you have `docker` and `docker-compose` installed and then build the server
images using

```shell
$ make && make up
```

This will fire up docker containers that can be used for trial and testing:

- One of docker containers runs OPA, using the policies in [`/docker/policy`](./docker/policy)

- The other two containers run the PAM modules. You can try running `sudo` and `SSH` on them.

#### SSH in

To SSH into the `frontend` container:

```bash
$ ssh -p 2222 ops@localhost -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null
```

Enter `ramesh` and `suresh` in the prompts that follow to get access.

This behavior can be modified by changing the policies in the directoy mentioned above.

#### Modify OpenSSH

To modify OpenSSH's behavior, edit [`/docker/etc/sshd_config`](./docker/etc/sshd_config).

For example you can change the line
```
AuthenticationMethods keyboard-interactive
```
to
```
AuthenticationMethods publickey,keyboard-interactive
```

The SSH server will now require both [this key](./docker/keys/id_rsa) and PAM module authorization before it grants access.

### Details

For a more details on how to install, run, debug the PAM module on your own machines, see [this README](./pam/README.md).