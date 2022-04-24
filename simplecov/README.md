# OPA Test Coverage to SimpleCov Format

Example Rego code to transform the OPA test coverage JSON report to
the [JSON representation](https://github.com/codeclimate-community/simplecov_json_formatter)
of the [SimpleCov](https://github.com/simplecov-ruby/simplecov) coverage report format.

## Why?

SimpleCov JSON is one of the supported report formats for
[AWS CodeBuild](https://docs.aws.amazon.com/codebuild/latest/userguide/build-spec-ref.html#reports-buildspec-file),
and likely other CI/CD tools as well.

## Example

Simply pipe the output of `opa test --coverage` into `opa eval` with the `simplecov.rego` file loaded:

```shell
$ opa test --coverage example \
| opa eval --format pretty \
           --stdin-input \
           --data simplecov.rego \
           data.simplecov.from_opa
```
**Output**
```json
{
  "coverage": {
    "example/policy.rego": {
      "lines": [
        null,
        null,
        1,
        1,
        null,
        null,
        0,
        0
      ]
    },
    "example/policy_test.rego": {
      "lines": [
        null,
        null,
        1,
        1
      ]
    }
  }
}
```

## Caveats

The OPA coverage report format does not report total number of lines in a file,
so the report will stop at the last occurence of a covered or not covered line.