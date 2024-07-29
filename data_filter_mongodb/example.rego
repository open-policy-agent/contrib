package details.authz

import rego.v1

allow if {
	input.method == "GET"
	some employee in data.employees
	employee.name == input.user
	employee.name == input.path[1]
}

allow if {
	input.method == "GET"
	some employee in data.employees
	employee.manager == input.user
	employee.name == input.path[1]
}

allow if {
	input.method == "GET"
	input.path = ["employees"]
	some employee in data.employees
	input.user == "danerys"
	employee.salary > 0
	employee.salary < 300000
}
