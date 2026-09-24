package aws.cloudformation

import rego.v1

# METADATA
# description: |
#   Get all resource types from the single-file JSON schema provided by AWS for the us-west-1
#   region. See https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-resource-specification.html
#   for a list of all available files in all regions
resource_types contains type if {
	resp := http.send({
		"method": "GET",
		"url": "https://d68hl49wbnanq.cloudfront.net/latest/gzip/CloudFormationResourceSpecification.json",
	})

	some type, _ in resp.body.ResourceTypes
}

# METADATA
# description: |
#   Transform all resource types into a their wildcard representation, i.e. "target".
#   Example: AWS::S3:Bucket => AWS::S3::*
resource_targets contains target if {
	some resource_type in resource_types
	target := concat("::", array.concat(array.slice(split(resource_type, "::"), 0, 2), ["*"]))
}

# METADATA
# description: |
#   Patch the provided input (the existing file) with the resource types provided at the AWS endpoint
output := object.union(
	input,
	{"handlers": {
		"preCreate": {"targetNames": resource_targets},
		"preUpdate": {"targetNames": resource_targets},
		"preDelete": {"targetNames": resource_targets},
	}},
)
