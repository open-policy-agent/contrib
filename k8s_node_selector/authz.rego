package system.authz

import rego.v1

# Deny access by default.
default allow := false

# Allow anonymous access to the default policy decision.
allow if {
	input.method == "POST"
	input.path == [""]
}

# Allow anonymouse access to /health otherwise K8s get 403 and kills pod.
allow if {
	input.path == ["health"]
}

# Allow authenticated traffic from the kube-mgmt sidecar
# in `./deploy.sh` the `{TOKEN_HERE}` variable is replaced

allow if {
	input.identity == "{TOKEN_HERE}"
}
