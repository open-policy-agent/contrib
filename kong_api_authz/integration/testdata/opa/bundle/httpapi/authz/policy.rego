package httpapi.authz

# Deny access by default
default allow = false

# Allows admin to access '/status' endpoints
allow {
  input.method == "GET"
  glob.match("/status**", ["/"], input.path)
  input.token.payload.role == "admin"
}