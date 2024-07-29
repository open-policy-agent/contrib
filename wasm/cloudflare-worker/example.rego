package example

import rego.v1

allow if {
	input.method == "GET"
}

allow if {
	input.method == "POST"
	input.path == "/api"
}
