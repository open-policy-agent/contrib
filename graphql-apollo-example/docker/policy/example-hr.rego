package graphqlapi.authz

# Allow HR members to get anyone's salary.
allowed_query(q) {
  selected_salary(q)
  input.user == hr[_]
}

# David is the only member of HR.
hr = [
  "david",
]
