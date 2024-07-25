package main

import rego.v1

find_vm_type(data_obj, service) := vm_type if {
	find_vm_type := [resource |
		some resource in data_obj.resources
		resource.identifier == service
	]

	vm_type := {
		"resource": find_vm_type,
		"present": count(find_vm_type) > 0,
	}
}
