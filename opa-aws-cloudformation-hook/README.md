# OPA AWS CloudFormation Hook

<p align="center">
    <img alt="OPA AWS CloudFormation Hook Diagram" src="docs/assets/opa-aws-cloudformation.svg">
</p>

This repository integrates AWS CloudFormation (CFN) with OPA using
[AWS Cloud Formation Hooks](https://aws.amazon.com/about-aws/whats-new/2022/02/aws-announces-general-availability-aws-cloudformation-hooks/).
Use this integration if you want to enforce policies over AWS resources (e.g., EC2 instances, S3 buckets, etc.)
provisioned with CloudFormation. For example, using this integration you can enforce policy across resources like:

* [EC2 Security Groups](./examples/policy/ec2/security_group/security_group.rego)
* [IAM Admin Rules](./examples/policy/iam/user/no_admin_test.rego)
* [S3 Public Access](./examples/policy/s3/bucket/public_access_test.rego)

> AWS Cloud Formation Hooks were added in February 2022. The feature is still relatively new for AWS Cloud Formation.
> If you run into any issues please report them [here](https://github.com/open-policy-agent/contrib/issues).

## How it Works

The OPA hook works by installing an
[AWS CloudFormation Hook](https://docs.aws.amazon.com/cloudformation-cli/latest/userguide/hooks-structure.html)
to your environment.

When creating, updating, or deleting a CloudFormation Stack, the hook is triggered to validate the configuration.
When used in conjunction with OPA, the hook will send the property information from each resource in a Stack to your
OPA server. When this information is received, OPA will validate the request against your defined policies and send
back any violations it may have found, which will stop the stack creation and log the violations to AWS CloudWatch.
If no violations are reported, the resources contained in the stack are created, updated or deleted accordingly.

**NOTE:** Installing OPA into your AWS environment is currently out of scope for this documentation. For local
development, a tool like [ngrok](https://ngrok.com/) could be used to point at an OPA running on your machine.

Want to try out this integration yourself? See the AWS Cloud Formation Hooks tutorial in the
[OPA documentation](https://www.openpolicyagent.org/docs/latest/aws-cloudformation-hooks/).

## Repository Contents

Provided in this repository, you'll find the code for the hook you'll deploy in your AWS account to enable OPA policy
enforcement for your CloudFormation resources under the `hooks` directory. See the
[OPA tutorial](https://www.openpolicyagent.org/docs/latest/aws-cloudformation-hooks/) on the topic for instructions on
how to quickly get started, or the
[development guide](https://docs.aws.amazon.com/cloudformation-cli/latest/userguide/hooks.html)
in AWS the documentation if you'd like to learn more about how it works.

To give you an idea about what policy for AWS CloudFormation Hooks might look like, this repository provides a number
of example resources and policies:

* The `examples/templates` directory contains example templates used for testing
* The `examples/policy` directory contains example policies

### Policy Development

In order to quickly iterate on changes in your Rego policies, you may use the `validate.py` tool provided under the
`test` directory. The tool allows you to test your policies against provided CloudFormation template files, without
actually submitting them to a hook installed in your environment. With an OPA server started with your policy
files loaded (e.g. `opa run --server --watch examples/policy`), you may use the tool like:

```shell
test/validate.py my-cloudformation-template.yaml
```

The tool will extract all resources found in the template and submit them to OPA one by one, in the same manner
the hook operates once installed. Should any violation be encountered, the tool will print them to the console.

### Deregistering the Hook

Deregistering a hook requires removal of not just the hook type, but also any versions of the hook deployed. In order
to help with that, you may use the `deregister-hook.sh` script provided in this repo, with the ARN of the hook provided
as the only argument:

```script
./deregister-hook.sh <ARN of your hook here>
```

## Community

For questions, discussions and announcements related to OPA, join the [OPA Slack](https://slack.openpolicyagent.org/).