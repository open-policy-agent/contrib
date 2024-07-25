# This package path should be passed with the authz_endpoint flag
# in the sshd PAM configuration file.
package sshd.authz

import rego.v1

default allow := false

allow if {
	# Verify user input.
	input.display_responses.last_name == "ramesh"
	input.display_responses.secret == "suresh"

	# Only allow running on host_id "frontend"
	input.pull_responses.files["/etc/host_identity.json"].host_id == "frontend"

	# Only authorize user "ops"
	input.sysinfo.pam_username = "ops"
}

errors contains "You cannot pass!" if {
	not allow
}
