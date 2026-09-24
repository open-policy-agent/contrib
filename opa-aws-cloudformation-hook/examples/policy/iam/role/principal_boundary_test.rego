package aws.iam.role_test

import rego.v1

import data.aws.iam.role.deny

import data.assertions.assert_empty

import data.test_helpers.create_with_properties
import data.test_helpers.with_properties

mock_create := create_with_properties("AWS::IAM::Role", "IAMRoleTest", {"AssumeRolePolicyDocument": {
	"Version": "2012-10-17",
	"Statement": [{
		"Action": "sts:AssumeRole",
		"Effect": "Allow",
		"Principal": {"Service": "codepipeline.amazonaws.com"},
	}],
}})

test_deny_auto_generated_name_not_excluded if {
	inp := object.union(mock_create, with_properties({
		"RoleName": "iam-not-excluded-cfn-hooks-cfn-stack-1-fail-046693375555",
		"PermissionsBoundary": "arn:aws:iam::555555555555:policy/invalid_s3_deny_permissions_boundary",
	}))

	# regal ignore:line-length
	deny["PermissionsBoundary arn:aws:iam::555555555555:policy/invalid_s3_deny_permissions_boundary is not allowed for IAMRoleTest"] with input as inp
}

test_deny_permission_boundary_not_set if {
	deny["PermissionsBoundary is not set for IAMRoleTest"] with input as mock_create
}

test_allow_permission_boundary_included if {
	inp := object.union(mock_create, with_properties({
		"RoleName": "cfn-hooks-pass-046693375555",
		"PermissionsBoundary": "arn:aws:iam::555555555555:policy/s3_deny_permissions_boundary",
	}))

	assert_empty(deny) with input as inp
}

test_allow_role_name_excluded if {
	inp := object.union(mock_create, with_properties({"RoleName": "excluded-cfn-hooks-stack1-046693375555"}))

	assert_empty(deny) with input as inp
}
