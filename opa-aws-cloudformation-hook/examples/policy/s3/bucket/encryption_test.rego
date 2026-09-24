package aws.s3.bucket_encryption_test

import rego.v1

import data.aws.s3.bucket.deny

import data.test_helpers.create_with_properties

test_deny_if_bucket_encryption_not_set if {
	inp := create_with_properties("AWS::S3::Bucket", "MyS3Bucket", {"BucketEncryption": {}})

	deny["bucket encryption is not enabled for bucket: MyS3Bucket"] with input as inp
}

test_deny_if_bucket_encryption_is_not_aws_kms if {
	inp := create_with_properties(
		"AWS::S3::Bucket", "MyS3Bucket",
		# regal ignore:line-length
		{"BucketEncryption": {"ServerSideEncryptionConfiguration": [{"ServerSideEncryptionByDefault": {"SSEAlgorithm": "aws:invalid"}}]}},
	)

	deny["encryption not set to aws:kms for bucket: MyS3Bucket"] with input as inp
}

test_allow_if_bucket_encryption_is_set if {
	inp := create_with_properties(
		"AWS::S3::Bucket",
		"MyS3Bucket",
		# regal ignore:line-length
		{"BucketEncryption": {"ServerSideEncryptionConfiguration": [{"ServerSideEncryptionByDefault": {"SSEAlgorithm": "aws:kms"}}]}},
	)

	not deny["bucket encryption is not enabled for bucket: MyS3Bucket"] with input as inp
	not deny["encryption not set to aws:kms for bucket: MyS3Bucket"] with input as inp
}
