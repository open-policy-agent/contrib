package pshelpers_test

import rego.v1

import data.pshelpers

test_limit_equals if {
	limit := 20
	set := [10, 12, 15, 20]

	not pshelpers.limit_exceeded(limit, set)
}

test_limit_is_exceeded if {
	limit := 20
	set := [10, 12, 15, 21]

	pshelpers.limit_exceeded(limit, set)
}

test_limit_not_exceeded if {
	limit := 20
	set := [10, 12, 15]

	not pshelpers.limit_exceeded(limit, set)
}

test_multiple_limits_exceeded if {
	limit := 20
	set := [21, 22, 1]

	pshelpers.limit_exceeded(limit, set)
}
