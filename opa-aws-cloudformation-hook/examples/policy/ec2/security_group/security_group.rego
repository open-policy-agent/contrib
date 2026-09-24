package aws.ec2.securitygroup

import rego.v1

deny contains msg if {
	input.resource.properties.SecurityGroupIngress[0].CidrIp == "0.0.0.0/0"

	msg := sprintf(
		"Security Group cannot contain rules allow all destinations (0.0.0.0/0 or ::/0): %s",
		[input.resource.id],
	)
}

deny contains msg if {
	input.resource.properties.SecurityGroupIngress[0].CidrIpv6 == "::/0"

	msg := sprintf(
		"Security Group cannot contain rules allow all destinations (0.0.0.0/0 or ::/0): %s",
		[input.resource.id],
	)
}
