## OPA test result to JUnit converter

Utility script to convert the output of OPA's (JSON) unit tests into the widely recognized JUnit XML format. This is primarly useful for CI/CD systems with ready made JUnit test report parsers, or any project where JUnit is already used for test reports and OPA tests should be incldued with those in a final test report.

## Run

Usage:

```sh
$ opa_test_to_junit.py <path>
```

Where `<path>` is a file containing the output of `opa test --format=json ...`

More conveniently the script also accepts input from stdin, allowing the output of the `opa test` command to be piped:

```sh
$ opa test --format=json <path> | opa_test_to_junit.py
```

Example output:

```xml
<?xml version='1.0' encoding='utf-8'?>
<testsuites errors="1" failures="6" tests="22" time="0.015">
  <testsuite errors="1" failures="6" hostname="localhost" name="data.kubernetes.authz" tests="22" time="0.015">
    <testcase classname="policy-test.rego" name="test_deny_by_default" time="0.001" />
    <testcase classname="policy-test.rego" name="test_allow_if_admin" time="0.002">
      <failure />
    </testcase>
    ...
```