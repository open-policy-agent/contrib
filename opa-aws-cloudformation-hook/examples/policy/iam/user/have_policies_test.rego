package aws.iam.user_have_policies_test

import rego.v1

import data.aws.iam.user.deny

import data.assertions.assert_in
import data.assertions.assert_not_in

import data.test_helpers.create_with_properties

mock_create := create_with_properties("AWS::IAM::User", "IAMUserTest", {
	"Policies": [{"PolicyDocument": {
		"PolicyName": "Test",
		"Statement": [{
			"Action": "'*'",
			"Effect": "Deny",
			"Resource": "'*'",
		}],
		"Version": "2012-10-17",
	}}],
	"UserName": "WithPolicy",
})

test_deny_policy_statement_undefined if {
	inp := create_with_properties("AWS::IAM::User", "IAMUserTest", {"AssumeRolePolicyDocument": {
		"Version": "2012-10-17",
		"Statement": [{
			"Action": "sts:AssumeRole",
			"Effect": "Allow",
			"Principal": {"Service": "codepipeline.amazonaws.com"},
		}],
	}})

	assert_in("IAM user does not have a policy statement: IAMUserTest", deny) with input as inp
}

test_allow_user_with_policy if {
	inp := create_with_properties("AWS::IAM::User", "IAMUserTest", {
		"Policies": [{"PolicyDocument": {
			"PolicyName": "Test",
			"Statement": [{
				"Action": "'*'",
				"Effect": "Deny",
				"Resource": "'*'",
			}],
			"Version": "2012-10-17",
		}}],
		"UserName": "WithPolicy",
	})

	assert_not_in("IAM user does not have a policy statement: IAMUserTest", deny) with input as inp
}

test_allow_user_with_managed_policy if {
	inp := create_with_properties(
		"AWS::IAM::User",
		"IAMUserTest",
		{"ManagedPolicyArns": ["arn:aws:iam::aws:policy/AWSDenyAll"]},
	)

	assert_not_in("IAM user does not have a policy statement: IAMUserTest", deny) with input as inp
}
