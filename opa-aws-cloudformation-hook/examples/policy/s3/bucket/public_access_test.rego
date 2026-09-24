package aws.s3.bucket_public_access_test

import rego.v1

import data.aws.s3.bucket.deny

import data.test_helpers.create_with_properties

test_deny_if_not_public_access_blocked if {
	inp := create_with_properties("AWS::S3::Bucket", "MyS3Bucket", {"PublicAccessBlockConfiguration": {
		"BlockPublicAcls": "false",
		"BlockPublicPolicy": "true",
		"IgnorePublicAcls": "true",
		"RestrictPublicBuckets": "false",
	}})

	deny["public access not blocked for bucket MyS3Bucket"] with input as inp
}

test_allow_if_public_access_blocked if {
	inp := create_with_properties("AWS::S3::Bucket", "MyS3Bucket", {"PublicAccessBlockConfiguration": {
		"BlockPublicAcls": "true",
		"BlockPublicPolicy": "true",
		"IgnorePublicAcls": "true",
		"RestrictPublicBuckets": "true",
	}})

	not deny["public access not blocked for bucket MyS3Bucket"] with input as inp
}

test_allow_if_excluded_prefix if {
	inp := create_with_properties("AWS::S3::Bucket", "MyS3Bucket", {"BucketName": "excluded-bucket"})

	not deny["public access not blocked for bucket MyS3Bucket"] with input as inp
}
