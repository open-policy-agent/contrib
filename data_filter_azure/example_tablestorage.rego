package example

import rego.v1

default allow := false

allow if {
	input.action in input.data[input.type][input.resourceName].actions
}

allow if {
	"*" in input.data[input.type][input.resourceName].actions
}

allow if {
	input.action in input.data[input.type]["*"].actions
}

allow if {
	"*" in input.data[input.type]["*"].actions
}
