# This package path should be passed with the pull_endpoint flag
# in the PAM configuration file.
package pull

import rego.v1

# JSON files to pull.
files := ["/etc/host_identity.json"]

# env_vars to pull.
env_vars := []
