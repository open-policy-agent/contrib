package httpapi.authz.hr

import rego.v1

# Allow HR members to get anyone's salary.
allow if {
	input.method == "GET"
	input.path = ["finance", "salary", _]
	input.user in members
}

# David is the only member of HR.
members := ["david"]
