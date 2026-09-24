package aws.s3.bucket

import rego.v1

deny contains msg if {
	not bucket_excluded_encryption
	not valid_bucket_encryption
	msg := sprintf("bucket encryption is not enabled for bucket: %s", [input.resource.id])
}

deny contains msg if {
	not bucket_excluded_encryption
	not valid_sse_algo
	msg := sprintf("encryption not set to aws:kms for bucket: %s", [input.resource.id])
}

valid_bucket_encryption if input.resource.properties.BucketEncryption != {}

valid_sse_algo if {
	input.resource.properties.BucketEncryption.ServerSideEncryptionConfiguration[0].ServerSideEncryptionByDefault.SSEAlgorithm == "aws:kms"
}

bucket_excluded_encryption if {
	some prefix in {"excluded-", "baseline-", "access-"}
	startswith(input.resource.properties.BucketName, prefix)
}
