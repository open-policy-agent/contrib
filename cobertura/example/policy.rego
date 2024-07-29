package policy

import rego.v1

deny contains "foo" if {
	input.foo == "bar"
}

deny contains "bar" if {
	input.bar == "foo"
}
