# Policy-driven SSH with OPA and a PAM module

This directory helps provide fine-grained, policy-based control over who can
login to each of your servers and containers.  It includes:

* A policy-enabled PAM module that you install on each of your servers or containers (/pam)
* Code showing you how to install and configure the PAM module and package servers as containers  (/docker)
* Tutorial demonstrating the PAM module works (this file)

## Demonstration

To demonstrate the benefits of the OPA PAM module, we created a docker image
with an ssh server and the PAM module installed and configured. The instructions
that follow take you through the steps of building the docker image, adding
policy, and sshing into the box.

What you'll see is that you can flexibly write and change policy that controls
who can ssh into the box and who can elevate to sudo rights without writing a
line of code, modifying LDAP groups, changing configuration management tools,
digging around inside the server's internals, or changing keys.  More
importantly, you use the same language for controlling ssh as you do for
controlling sudo, which is the same language you use for controlling Kubernetes,
docker, and any of the other systems integrated with OPA.

### Setup

To get started, make sure you have docker installed and then build the server
images using `make`. After this completes, `docker images` will show the
following images: `ssh-backend`, `ssh-webapp`, and `pam- builder`. The builder
image is used for compiling the PAM module, and the other two images are our
example "servers" with an ssh server installed with the PAM module.


### SSH and sudo access policy

Our desired policy for ssh authorization is the following:

* Admin users can ssh into any server
* Developers can ssh into servers whose role match the developers role

To implement this policy, we need two data sources: one that specifies the users
in the admin group and another describing the role for each user. For this
example, these datasets are specified in the files
[docker/policy/users.json](docker/policy/users.json) and
[docker/policy/contribuers.json](docker/policy/contribuers.json).  (In reality,
this information would be pushed into OPA using its API.)

The desired sudo authorization policy is:

* Admin users can use sudo
* No other uses can use sudo

Again, this policy relies on external data --- a list of admin users.

### Running the Servers

Time to put the policy into action. The two docker-containers above call into
OPA for policy decisions, so for the demo we have 3 containers (2 "servers" and
1 OPA).  To simplify startup, we included a docker-compose file. Make sure that
you have [docker-compose](https://docs.docker.com/compose/install/) installed,
and then use the command:

```
make up
```

Using the docker-compose file, OPA is preconfigured with policies for ssh and
sudo authorization. When a user attempts to connect to the ssh server,
authentication is performed first. For our demo, we use standard Linux users and
ssh keys for this. Once a user is authenticated, they are authorized by the PAM
module. To see this in action, we pre-configured the server containers with a
few users and group those users into two categories: engineering or admin.  The
key you will use for sshing can be found in `docker/id_rsa` (which you may need
to set permissions for with `chmod 0600 docker/id_rsa`).

### Walkthrough

Because this demo has you sshing into docker containers, you will ssh into the
IP address that docker is using (aliased as `docker_ip` shown below), and you
will use port number 2222 for the web-app container and 2223 for the backend
container. Usually the servers/containers you will install the PAM module onto
will have public IPs or hostnames, and you'll ssh into them in the usual way,
providing just the IP or hostname.

For convenience, let's set up an alias to the docker IP as shown below.

```
# The following line is only for Mac users using docker machine
$ docker_ip=`docker-machine ip default`
# Linux users and Mac users not using docker machine should use the following
$ docker_ip=localhost
```

Now let's try and access the servers using ssh. Recall that users in the admin
group (as defined in policy) have ssh access to any server and can perform sudo
commands. The user `ops-dev` is an admin, so run the following commands and see
that they all succeed.

```
# ssh into webapp, with sudo
$ ssh -i docker/id_rsa -p 2222 ops-dev@$docker_ip
$ sudo ls /
$ logout
# ssh into backend, with sudo
$ ssh -i docker/id_rsa -p 2223 ops-dev@$docker_ip
$ sudo ls /
$ logout
```

Let's try a user not in the admin group. Recall that a non-admin can ssh into
any server running an app they wrote code for.   The user `web-dev` contributes
to the webapp server, but not the backend. `web-dev` can ssh into servers with
the `webapp` role, but not any other servers. Further, this user should not be
able to use sudo.

```
# ssh into webapp, but no sudo
$ ssh -i docker/id_rsa -p 2222 web-dev@$docker_ip
$ sudo ls /
$ logout
# no ssh into backend
$ ssh -i docker/id_rsa -p 2223 web-dev@$docker_ip
```

To this point, the policy could have been hard coded into the PAM module. OPA,
however, allows us to dynamically modify the policies and the data that policies
rely on. Consider the user "Pam". In our initial data, Pam has contributed to
the WebApp, so can only ssh into the container with the webapp role. Suppose
that Pam makes changes to the backend as well, and so Pam should be added to the
list of Backend contributers. We use OPA's API to update the data like so:

```
curl -X PUT $docker_ip:8181/v1/data/io/vcs/contributers -d \
'{
    "WebApp": {
        "Contributers": ["web-dev", "Pam", "Hans"]
    },
    "Backend": {
        "Contributers": ["Pam", "backend-dev", "Stan", "Hans"]
    }
}'
```

Now `ssh -i docker/id_rsa -p 2223 web-dev@$docker_ip` is allowed.

As another example, suppose that a bug occurs in production and Pam needs sudo
access to debug. We can update the list of admins to give Pam sudo access like
so:

```
curl -X PUT $docker_ip:8181/v1/data/io/directory/users -d \
'{
    "engineering": ["web-dev", "backend-dev" ,"ops-dev", "Stan", "Hans", "Pam", "Sam", "Jan"],
    "admin": ["Pam", "ops-dev"]
}'
```

Now Pam has ssh and sudo access on all servers.


### Orchestration of Multiple Servers

The previous example shows how OPA and PAM can be used for authorization on a
single server, but in practice policy needs to be enforced across many servers.
Because OPA is a host-local daemon, we recommend using a standalone, centralized
service (like etcd) to store policy and data centrally and a side-car (or
wrapper) for OPA that keeps all the OPAs up to date with that central service.
(OPA was purpose-built to be a host-local daemon to ensure high-availability and
high-performance even in the presence of network partitions and leaves the
problem of policy/data replication up to the environment in which it is
deployed.)



