package system

import rego.v1

# regal ignore:pointless-reassignment
main := allow

default allow := false

allow if {
	input.method = "GET"
	input.path = [""]
}

allow if {
	input.method = "GET"
	input.path = ["cars"]
}

allow if {
	input.method = "GET"
	input.path = ["cars", car_id]
}

allow if {
	input.method = "PUT"
	input.path = ["cars", car_id]
	has_role("manager")
}

allow if {
	input.method = "DELETE"
	input.path = ["cars", car_id]
	has_role("manager")
}

allow if {
	input.method = "GET"
	input.path = ["cars", car_id, "status"]
	employees[input.user]
}

allow if {
	input.method = "PUT"
	input.path = ["cars", car_id, "status"]
	has_role("car_admin")
}

has_role(name) if {
	employee := employees[input.user] # regal ignore:external-reference
	employee.roles[name]
}

employees := {
	"alice": {"roles": {"manager", "car_admin"}},
	"james": {"roles": {"manager"}},
	"kelly": {"roles": {"car_admin"}},
}
