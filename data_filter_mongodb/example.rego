package details.authz

allow {
  input.method == "GET"
  employee := data.employees[_]
  employee.name == input.user
  employee.name == input.path[1]
}

allow {
  input.method == "GET"
  employee := data.employees[_]
  employee.manager == input.user
  employee.name == input.path[1]
}

allow {
  input.method == "GET"
  input.path = ["employees"]
  employee := data.employees[_]
  input.user == "danerys"
  employee.salary > 0
  employee.salary < 300000
}

