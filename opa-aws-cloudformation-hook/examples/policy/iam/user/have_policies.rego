package aws.iam.user

import rego.v1

deny contains msg if {
	not user_policies_exist

	msg := sprintf("IAM user does not have a policy statement: %s", [input.resource.id])
}

user_policies_exist if input.resource.properties.Policies

user_policies_exist if input.resource.properties.ManagedPolicyArns
