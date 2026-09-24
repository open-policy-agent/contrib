package aws.ec2.securitygroup_test

import rego.v1

import data.aws.ec2.securitygroup.deny

import data.assertions.assert_empty

import data.test_helpers.create_with_properties

test_deny_if_security_group_allows_all_destinations if {
	inp := create_with_properties("AWS::EC2::SecurityGroup", "SecurityGroup", {"SecurityGroupIngress": [{
		"CidrIp": "0.0.0.0/0",
		"IpProtocol": "-1",
	}]})

	deny["Security Group cannot contain rules allow all destinations (0.0.0.0/0 or ::/0): SecurityGroup"] with input as inp
}

test_allow_if_security_group_cidr_is_set if {
	inp := create_with_properties("AWS::EC2::SecurityGroup", "SecurityGroup", {"SecurityGroupIngress": [{
		"CidrIp": "10.0.0.0/16",
		"IpProtocol": "-1",
	}]})

	assert_empty(deny) with input as inp
}
