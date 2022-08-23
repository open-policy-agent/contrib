package httpapi.authz

# Allow HR members to get anyone's salary.
allow {
  input.method == "GET"
  input.path = ["finance", "salary", _]
  input.user == hr[_]
}

# David is the only member of HR.
hr = [
  "david",
]
