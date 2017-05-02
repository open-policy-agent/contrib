# HTTP API Authorization with OPA

This directory helps provide fine-grained, policy-based control over who
can run which RESTful APIs.  It includes

* A sample web application that asks OPA for authorization before executing an API call
* A default policy that allows `/finance/salary/<user>` for `<user>` and for `<user>`'s manager


## Setup:

The web application and OPA both run in docker-containers.  For convenience we
included a docker-compose file, so you'll want
[docker-compose](https://docs.docker.com/compose/install/) installed.

To build the containers and get them started, use the following make commands.

```
make       # build the containers with docker
make up    # start the containers with docker-compose
```

## Walkthrough

You'll be sending HTTP requests to a docker container, so it's useful to
setup an alias to the IP that docker uses.

```
# The following line is only for Mac users using docker machine
$ docker_ip=`docker-machine ip default`
# Linux users and Mac users not using docker machine should use the following
$ docker_ip=localhost
```

In a different terminal window from where you ran `make up`,
run a `curl` command to check that `alice` can see her
own salary.

```
$ docker_ip=localhost
$ curl --user alice:password $docker_ip:5000/finance/salary/alice
```

Also check that `bob` can see `alice`'s salary (because `bob` is
`alice`'s manager).

```
$ docker_ip=localhost
$ curl --user bob:password $docker_ip:5000/finance/salary/alice
```

But notice that `bob` cannot see `charlie`'s salary (because
`bob` is not `charlie`'s manager).

```
$ curl --user bob:password $docker_ip:5000/finance/salary/alice
```

## Change the policy

If you want to change the policy, go into `docker/policies` and
change the files you see there.  You can change both the high-level
rules and who is whose manager.  Then Ctrl+C in the terminal window
where you ran `make up` and run `make up` again.


