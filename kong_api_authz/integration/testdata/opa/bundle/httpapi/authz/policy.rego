package httpapi.authz

import rego.v1

# Deny access by default
default allow := false

# Allows admin to access '/status' endpoints
allow if {
	input.method == "GET"
	glob.match("/status**", ["/"], input.path)
	input.token.payload.role == "admin"
}
