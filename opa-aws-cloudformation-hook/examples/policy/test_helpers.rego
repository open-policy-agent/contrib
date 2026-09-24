package test_helpers

import rego.v1

create_with_properties(type, id, properties) := object.union(mock_create(type, id), with_properties(properties))

mock_create(type, id) := {
	"action": "CREATE",
	"hook": "OpenPolicyAgent::OPA::Hook",
	"resource": {
		"id": id,
		"name": type,
		"properties": {},
		"type": type,
	},
}

with_properties(obj) := {"resource": {"properties": obj}}
