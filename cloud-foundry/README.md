# opa-policy
This directory provides configuration policies and validators for you Cloud Foundry Deployments. \
With the use of [OPA](https://www.openpolicyagent.org/), bundles of policies to be used with [conftest](https://github.com/instrumenta/conftest) for fitness function testing your platform.

**OPA** provides the policy rules as well as policy language `rego`. **Conftest** provides developers with local tooling to do static analysis of configuration files so that engineers can shift left on testing infrastructure as code.

Why you should care? Being able to test baseline values before attempting to deploy allows operators to *fail fast*. 
These policies along with these tools are used to help bridge that gap.

# Using this Repository

## Examples

This are intended to be example tests that can continue to be contributed to. The both validate values and check configuration.

## Testing Configuration and Validation

The main point of this repo is to run [conftest](https://github.com/instrumenta/conftest) against the OPA policies contained within. 

For example, suppose you're a developer writing a configuration file for the Pivotal Application Service or a Cloud Foundry Application Runtime. In the root of this directory do the following.

```sh
$ cat > /tmp/pas.yml <<-EOF
---
product-properties:
  ".properties.credhub_hsm_provider_partition_password":
    value:
    - primary: false
  ".properties.credhub_key_encryption_passwords":
    value:
    - primary: '1234567890123456789'
EOF
```

Let's suppose you want to do a sanity check on whether you've made a mistake in your YAML file. You can use conftest to run OPA policies that check your YAML files for obvious mistakes.

A few notes for below. `-p` indicates the policy being run. Any arguments without flags preceding are treated as the files being evaluated. In this case `/tmp/pas.yml` is being evaluated.

```sh
-> % conftest test --namespace credhub -p cloud-foundry/polices/pas/credhub-key/credhub_key.rego /tmp/pas.yml
PASS - /tmp/pas.yml - data.credhub.deny_if_not_exactly_one_primary
PASS - /tmp/pas.yml - data.credhub.deny_not_enough_chars
```

Check out the OPA and conftest communities for information on running them. This grouping of policies is intended to be a starting point for cloud foundry users.

## .rego tests
If you write code you should be testing it. Releases to this project only accept code that has passed tests with the ci tool. Please test your rego if you plan to contribute

The policies are separated by packages so run the following command to capture all of them.

```shell
opa test -vl policies/*
```

You can also test on commit by using the included githook. From root of this directory run

```shell
cp ./.githooks/* ./.git/hooks/
```
