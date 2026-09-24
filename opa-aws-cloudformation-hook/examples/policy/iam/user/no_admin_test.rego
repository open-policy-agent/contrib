package aws.iam.user_no_admin_test

import rego.v1

import data.aws.iam.user.deny

import data.assertions.assert_empty
import data.assertions.assert_in
import data.assertions.assert_not_in

import data.test_helpers.create_with_properties
import data.test_helpers.with_properties

mock_create := create_with_properties("AWS::IAM::User", "IAMUserTest", {
	"AssumeRolePolicyDocument": {
		"Version": "2012-10-17",
		"Statement": [{
			"Action": "sts:AssumeRole",
			"Effect": "Allow",
			"Principal": {"Service": "codepipeline.amazonaws.com"},
		}],
	},
	"PermissionsBoundary": "arn:aws:iam::555555555555:policy/s3_deny_permissions_boundary",
})

test_deny_policy_statement_undefined if {
	inp := mock_create

	assert_in("please limit the scope for IAM user: IAMUserTest", deny) with input as inp
}

test_deny_user_with_admin if {
	inp := object.union(mock_create, with_properties({
		"Policies": [
			{"PolicyDocument": {
				"PolicyName": "Test1",
				"Statement": [{
					"Action": "'Create'",
					"Effect": "Allow",
					"Resource": "'*'",
				}],
				"Version": "2012-10-17",
			}},
			{"PolicyDocument": {
				"PolicyName": "Test2",
				"Statement": [{
					"Action": "'*'",
					"Effect": "Allow",
					"Resource": "'*'",
				}],
				"Version": "2012-10-17",
			}},
		],
		"UserName": "WithInlineAdminPolicy",
	}))

	msg := "please limit the scope for IAM user: IAMUserTest"

	assert_in(msg, deny) with input as inp
}

test_allow_user_without_admin if {
	inp := object.union(mock_create, with_properties({
		"Policies": [{"PolicyDocument": {
			"PolicyName": "Test1",
			"Statement": [{
				"Action": "'Create'",
				"Effect": "Allow",
				"Resource": "'*'",
			}],
			"Version": "2012-10-17",
		}}],
		"UserName": "WithInlineAdminPolicy",
	}))

	msg := "please limit the scope for IAM user: IAMUserTest"

	assert_not_in(msg, deny) with input as inp
}

test_allow_user_limited_scope if {
	inp := object.union(mock_create, with_properties({
		"Policies": [{"PolicyDocument": {
			"PolicyName": "Test",
			"Statement": [{
				"Action": "'ec2:*'",
				"Effect": "Allow",
				"Resource": "'*'",
			}],
			"Version": "2012-10-17",
		}}],
		"UserName": "WithInlineEc2Policy",
	}))

	assert_empty(deny) with input as inp
}
