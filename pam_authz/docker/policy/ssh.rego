package ssh.authz

# ssh policy goal: give access to administrators and individuals who have
# contributed code to an app running on the server

import input.user
import input.host_identity

import data.io.directory.users
import data.io.vcs.contributors

# By default, users are not authorized
default allow = false

# If a group is specified, then members of that group are always authorized
allow {
    user = users["admin"][_]
}

# Authorize users who contributed to the app running on this server
# Users who match this rule or the rule above will be given access
allow {
    user = contributors[host_identity.host_id].Contributors[_]
}

# If the user is not authorized, then include an error message in the response
errors["Request denied by administrative policy"] {
  not allow
}
