package aws.iam.user

import rego.v1

deny contains msg if {
	not managed_policy_exist
	not valid_iam_scope

	msg := sprintf("please limit the scope for IAM user: %s", [input.resource.id])
}

valid_iam_scope if every policy in input.resource.properties.Policies {
	every statement in policy.PolicyDocument.Statement {
		statement.Action != "'*'"
	}
}

managed_policy_exist if input.resource.properties.ManagedPolicyArns
