# This package path should be passed with the authz_endpoint flag
# in the PAM configuration file.
package common.authz

import input.display_responses
import input.pull_responses
import input.sysinfo

default allow = false

allow {
	# Verify user input.
	display_responses.last_name = "ramesh"
	display_responses.secret = "suresh"

	# Only allow running on host_id "frontend"
	pull_responses.files["/etc/host_identity.json"].host_id = "frontend"

	# Only authorize user "ops"
	sysinfo.pam_username = "ops"
}

errors["You cannot pass!"] {
	not allow
}