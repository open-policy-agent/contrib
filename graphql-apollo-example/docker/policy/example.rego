package graphqlapi.authz

import future.keywords.in
import future.keywords.every

subordinates = {"alice": [], "charlie": [], "bob": ["alice"], "betty": ["charlie"]}

query_ast := graphql.parse(input.query, input.schema)[0] # If validation fails, the rules depending on this will be undefined.

default allow := false

allow {
    employeeByIDQueries != {}
    every query in employeeByIDQueries {
      allowed_query(query)
    }
}

# Allow users to see the salaries of their subordinates. (variable case)
allowed_query(q) {
    selected_salary(q)
    varname := variable_arg(q, "id")
    input.variables[varname] in subordinates[input.user] # Do value lookup from the 'variables' object.
}
# Allow users to see the salaries of their subordinates. (constant value case)
allowed_query(q) {
    selected_salary(q)
    username := constant_string_arg(q, "id")
    username in subordinates[input.user]
}

# Helper rules.

# Allow users to get their own salaries. (variable case)
allowed_query(q) {
    selected_salary(q)
    varname := variable_arg(q, "id")
    input.user == input.variables[varname] # Do value lookup from the 'variables' object.
}

# Allow users to get their own salaries. (constant value case)
allowed_query(q) {
    selected_salary(q)
    username := constant_string_arg(q, "id")
    input.user == username
}


# Helper functions.

# Build up an object with all queries of interest as values.
employeeByIDQueries[value] {
    some value
	walk(query_ast, [_, value])
    value.Name == "employeeByID"
    count(value.SelectionSet) > 0 # Ensure we latch onto an employeeByID query.
}

# Extract the string value of a constant value argument.
constant_string_arg(value, argname) := arg.Value.Raw {
    some arg in value.Arguments
    arg.Name == argname
    arg.Value.Kind == 3
}

# Extract the variable name for a variable argument.
variable_arg(value, argname) := arg.Value.Raw {
    some arg in value.Arguments
    arg.Name == argname
    arg.Value.Kind == 0
}

# Ensure we're dealing with a selection set that includes the "salary" field.
selected_salary(value) := value.SelectionSet[_].Name == "salary"

