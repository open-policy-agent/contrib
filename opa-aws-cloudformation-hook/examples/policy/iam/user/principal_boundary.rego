package aws.iam.user

import rego.v1

excluded_principal_prefixes := ["excluded", "iam-excluded", "test-excluded-user"]

deny contains msg if {
	not excluded_principal_name
	not permission_boundary_exists

	msg := sprintf("PermissionsBoundary is not set for %s", [input.resource.id])
}

deny contains msg if {
	not excluded_principal_name
	permission_boundary_exists
	not valid_permission_boundary

	msg := sprintf(
		"PermissionsBoundary %s is not allowed for %s",
		[input.resource.properties.PermissionsBoundary, input.resource.id],
	)
}

excluded_principal_name if {
	name := input.resource.properties.UserName
	some prefix in excluded_principal_prefixes
	startswith(name, prefix)
}

permission_boundary_exists if input.resource.properties.PermissionsBoundary

valid_permission_boundary if {
	input.resource.properties.PermissionsBoundary == "arn:aws:iam::555555555555:policy/s3_deny_permissions_boundary"
}
