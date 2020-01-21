package system.authz

# Deny access by default.
default allow = false

# Allow anonymous access to the default policy decision.
allow {
    input.method = "POST"
    input.path = [""]
}


# Allow anonymouse access to /health otherwise K8s get 403 and kills pod. 
allow {
    input.path = ["health"]
}

# Allow authenticated traffic from the kube-mgmt sidecar
# in `./deploy.sh` the `{TOKEN_HERE}` variable is replaced
allow {
  "{TOKEN_HERE}" = input.identity
}