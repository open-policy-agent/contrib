package puppet.authz

import rego.v1

# regal ignore:unresolved-import
import data.git

default allow := false

allow if not deny

deny if {
	some resource_index, resource in input.puppet.catalog.resources
	resource.type == "File"
	startswith(resource.title, "/etc/app")
	email := resource_author[resource_index]
	not email in app_team
}

deny if {
	some resource_index, resource in input.puppet.catalog.resources
	resource.type == "File"
	startswith(resource.title, "/etc/infra")
	email := resource_author[resource_index]
	not email in infra_team
}

resource_author[resource_index] := email if {
	# For each "File" resource...
	some resource_index, resource in input.puppet.catalog.resources
	resource.type == "File"

	# Compute the Puppet manifest filename relative to Git repository root...
	prefix_length := count(source_root_dir)
	local_file := substring(resource.file, prefix_length, -1)

	# Lookup author for resource using the line number and Git blame data.
	blame_entry := git[local_file]
	email := blame_entry[resource.line].Author
}

app_team := {
	"alice@acmecorp.com",
	"bob@acmecorp.com",
}

infra_team := {
	"betty@acmecorp.com",
	"chris@acmecorp.com",
}

source_root_dir := "/code/"
