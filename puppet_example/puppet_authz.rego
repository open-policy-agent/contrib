package puppet.authz

import data.git
import input.puppet.catalog

default allow = false

allow { not deny }

deny {
    resource = catalog.resources[resource_index]
    resource.type = "File"
    startswith(resource.title, "/etc/app")
    resource_author[resource_index] = email
    not app_team[email]
}

deny {
    resource = catalog.resources[resource_index]
    resource.type = "File"
    startswith(resource.title, "/etc/infra")
    resource_author[resource_index] = email
    not infra_team[email]
}

resource_author[resource_index] = email {

    # For each "File" resource...
	resource = catalog.resources[resource_index]
    resource.type = "File"

    # Compute the Puppet manifest filename relative to Git repository root...
	count(source_root_dir, prefix_length)
    substring(resource.file, prefix_length, -1, local_file)

    # Lookup author for resource using the line number and Git blame data.
    blame_entry = git[local_file]
    email = blame_entry[resource.line]["Author"]
}

app_team = {
    "alice@acmecorp.com",
    "bob@acmecorp.com",
}

infra_team = {
    "betty@acmecorp.com",
    "chris@acmecorp.com",
}

source_root_dir = "/code/"
