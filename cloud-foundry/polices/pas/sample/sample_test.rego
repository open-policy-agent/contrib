package main

test_deny_suite {
  deny[_] with input as {"path": "blah", "method": "POST"}
  deny[_] with input as {"path": "nope", "method": "POST"}
}

