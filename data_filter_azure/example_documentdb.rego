package documentdb.example

import rego.v1

default allow := false

allow if {
	count(allowed) > 0
}

allowed contains user if {
	some user, _ in data.permissions
	user.registry == input.registry
	user.user == input.user
	type_in_input
	matched_name
	contains(user.map[_].actions, input.action)
}

allowed contains user if {
	some user, _ in data.permissions
	user.registry == input.registry
	user.user == input.user
	type_in_input
	wildcard_name
	contains(user.map[_].actions, input.action)
}

allowed contains user if {
	some user, _ in data.permissions
	user.registry == input.registry
	user.user == input.user
	type_in_input
	matched_name
	contains(user.map[_].actions, "*")
}

allowed contains user if {
	some user, _ in data.permissions
	user.registry == input.registry
	user.user == input.user
	type_in_input
	wildcard_name
	contains(user.map[_].actions, "*")
}

type_in_input if {
	some user, _ in data.permissions
	some map in user.map
	map.type == input.type
}

wildcard_name if {
	some user, _ in data.permissions
	some map in user.map
	map.name == "*"
}

matched_name if {
	some user, _ in data.permissions
	some map in user.map
	map.name == input.resourceName
}
