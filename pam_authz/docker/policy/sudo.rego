package sudo.authz

# sudo policy goal: only allow administrators to have access

import input.user
import input.host_identity

import data.io.directory.users
import data.io.vcs.contributors

default allow = false

# Only members of the specified group are authorized to use sudo
allow {
    user = users["admin"][_]
}

# If the user is not authorized, then include an error message in the response and include
# the group the user must be part of as part of the response
errors["Request denied by administrative policy, user is not a member of the group admin"] {
  not allow
}
