package aws.s3.bucket

import rego.v1

deny contains msg if {
	not bucket_excluded_public
	not public_access_blocked

	msg := sprintf("public access not blocked for bucket %s", [input.resource.id])
}

bucket_excluded_public if {
	some prefix in {"excluded-", "baseline-", "secure-"}
	startswith(input.resource.properties.BucketName, prefix)
}

public_access_blocked if {
	every property in ["BlockPublicAcls", "BlockPublicPolicy", "IgnorePublicAcls", "RestrictPublicBuckets"] {
		input.resource.properties.PublicAccessBlockConfiguration[property] == "true"
	}
}
