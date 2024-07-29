package policy_test

import rego.v1

import data.policy

test_deny if {
	count(policy.deny) == 1 with input as {"foo": "bar"}
}
