package aws.eks.cluster

import rego.v1

deny contains msg if {
	not public_endpoint_disabled
	msg := sprintf("public endpoint needs to be disabled for cluster: %s", [input.resource.id])
}

deny contains msg if {
	public_endpoint_disabled
	not private_endpoint_enabled
	msg := sprintf("invalid configuration, please enable private api for cluster: %s", [input.resource.id])
}

public_endpoint_disabled if input.resource.properties.ResourcesVpcConfig.EndpointPublicAccess == "false"

private_endpoint_enabled if input.resource.properties.ResourcesVpcConfig.EndpointPrivateAccess == "true"
