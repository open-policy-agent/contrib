# OPA Test Coverage to Cobertura Format

Example code to transform the OPA test coverage JSON report to
the [Cobertura](https://github.com/cobertura/cobertura/blob/master/cobertura/src/test/resources/dtds/coverage-04.dtd) coverage report format.

## Why?

Cobertura is one of the supported report formats for
[Jenkins Coverage API plugin](https://github.com/jenkinsci/code-coverage-api-plugin), which is widely used.

## Example

Simply pipe the output of `opa test --coverage` into `opa eval` with the `simplecov.rego` file loaded:

1. generate json format coverage result with `--coverage` and pipe it into a file
2. call `python opa_coverage_to_cobertura.py <input json> <output xml>` to generate xml format coverage report
```shell
$ opa test --coverage example > coverage.json
$ python opa_coverage_to_cobertura.py coverage.json coverage.xml
```
**Output**
```xml
<?xml version='1.0' encoding='utf-8'?>
<coverage lines-covered="4" line-rate="0.6665000000000001" lines-valid="6" complexity="0" version="0.1" timestamp="1683450316053">
  <packages>
    <package complexity="0" line-rate="0.6665000000000001" name="">
      <classes>
        <class complexity="0" line-rate="0.5" filename="example/policy.rego" name="example/policy.rego">
          <methods />
          <lines>
            <line number="3" hits="1" />
            <line number="4" hits="1" />
            <line number="7" hits="0" />
          </lines>
        </class>
        <class complexity="0" line-rate="1.0" filename="example/policy_test.rego" name="example/policy_test.rego">
          <methods />
          <lines>
            <line number="3" hits="1" />
          </lines>
        </class>
      </classes>
    </package>
  </packages>
</coverage>
```
