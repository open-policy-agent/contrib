# Policy-driven SSH and sudo with OPA and a PAM module

This directory helps provide fine-grained, policy-based control over who can
ssh and sudo into each of your servers and containers.

You can find a step-by-step tutorial at [openpolicyagent.org/tutorials/ssh-sudo-authorization/](http://www.openpolicyagent.org/tutorials/ssh-sudo-authorization/).


## Directory contents

This directory includes:

* A policy-enabled PAM module that you install on each of your servers or containers (/pam)
* Code showing you how to install and configure the PAM module and package servers as containers  (/docker)

### Setup

To get started, make sure you have docker installed and then build the server
images using

```shell
$ make
```
