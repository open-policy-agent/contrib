package policy_test

import rego.v1

import data.policy

test_deny if {
	policy.deny with input as {"foo": "bar"}
}
