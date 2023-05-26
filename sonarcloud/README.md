# OPA Test Coverage to SonarCloud Format

Example code to transform the OPA test coverage JSON report to
the [SonarCloud](https://docs.sonarqube.org/latest/analyzing-source-code/test-coverage/generic-test-data/) coverage report format.

## Why?

SonarCloud allows collecting test coverage data to generate rich reports. Integration can be done using [generic test data format](https://docs.sonarqube.org/latest/analyzing-source-code/test-coverage/generic-test-data/).

## Example
1. generate json format coverage result with `--coverage` and pipe it into a file
2. call `python opa_coverage_to_sonarcloud.py <input json> <output xml>` to generate xml format coverage report
```shell
$ opa test --coverage example > coverage.json
$ python opa_coverage_to_sonarcloud.py coverage.json coverage.xml
```
**Output**
```xml
<?xml version='1.0' encoding='utf-8'?>
<coverage version="1">
  <file path="policy.rego">
    <lineToCover lineNumber="3" covered="true" />
    <lineToCover lineNumber="4" covered="true" />
    <lineToCover lineNumber="7" covered="false" />
    <lineToCover lineNumber="8" covered="false" />
  </file>
  <file path="policy_test.rego">
    <lineToCover lineNumber="3" covered="true" />
    <lineToCover lineNumber="4" covered="true" />
  </file>
</coverage>
```
