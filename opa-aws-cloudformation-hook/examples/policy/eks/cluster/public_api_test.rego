package aws.eks.public_api_test

import rego.v1

import data.aws.eks.cluster.deny

import data.assertions.assert_in
import data.assertions.assert_not_in

import data.test_helpers.create_with_properties

test_allow_cluster_private_api if {
	inp := create_with_properties("AWS::EKS::Cluster", "EksCluster", {"ResourcesVpcConfig": {
		"RoleArn": "<MY_EKS_SERVICE_ROLE_ARN>",
		"SubnetIds": ["<MY_SUBNET_ID>"],
		"EndpointPublicAccess": "false",
		"EndpointPrivateAccess": "true",
	}})

	msg := "public endpoint needs to be disabled for cluster: EksCluster"
	assert_not_in(msg, deny) with input as inp
}

test_deny_cluster_public_api if {
	inp := create_with_properties("AWS::EKS::Cluster", "EksCluster", {"ResourcesVpcConfig": {
		"RoleArn": "<MY_EKS_SERVICE_ROLE_ARN>",
		"SubnetIds": ["<MY_SUBNET_ID>"],
		"EndpointPublicAccess": "true",
		"EndpointPrivateAccess": "true",
	}})

	msg := "public endpoint needs to be disabled for cluster: EksCluster"
	assert_in(msg, deny) with input as inp
}

test_deny_cluster_no_access if {
	inp := create_with_properties("AWS::EKS::Cluster", "EksCluster", {"ResourcesVpcConfig": {
		"RoleArn": "<MY_EKS_SERVICE_ROLE_ARN>",
		"SubnetIds": ["<MY_SUBNET_ID>"],
		"EndpointPublicAccess": "false",
		"EndpointPrivateAccess": "false",
	}})

	msg := "invalid configuration, please enable private api for cluster: EksCluster"
	assert_in(msg, deny) with input as inp
}
