package aws.s3.bucket_logging_test

import rego.v1

import data.aws.s3.bucket.deny

import data.test_helpers.create_with_properties

test_deny_if_logging_prefix_invalid if {
	inp := create_with_properties(
		"AWS::S3::Bucket", "MyS3Bucket",
		{"LoggingConfiguration": {"LogFilePrefix": "invalid"}},
	)

	deny["logging prefix is not set correctly for bucket: MyS3Bucket"] with input as inp
}

test_deny_if_logging_bucket_invalid if {
	inp := create_with_properties(
		"AWS::S3::Bucket", "MyS3Bucket",
		{"LoggingConfiguration": {"DestinationBucketName": "invalid"}},
	)

	deny["logging destination bucket is not correct: MyS3Bucket"] with input as inp
}

test_deny_if_logging_configuration_unset if {
	mock_create := create_with_properties("AWS::S3::Bucket", "MyS3Bucket", {})

	deny["logging destination bucket is not correct: MyS3Bucket"] with input as mock_create
	deny["logging prefix is not set correctly for bucket: MyS3Bucket"] with input as mock_create
}

test_allow_if_prefix_and_destination_set if {
	inp := create_with_properties("AWS::S3::Bucket", "MyS3Bucket", {"LoggingConfiguration": {
		"LogFilePrefix": "s3-logs",
		"DestinationBucketName": "my-logging-bucket",
	}})

	not deny["logging destination bucket is not correct: MyS3Bucket"] with input as inp
	not deny["logging prefix is not set correctly for bucket: MyS3Bucket"] with input as inp
}

test_allow_if_bucket_name_set if {
	inp := create_with_properties("AWS::S3::Bucket", "MyS3Bucket", {
		"BucketName": "My-Bucket",
		"LoggingConfiguration": {
			"LogFilePrefix": "s3-logs-My-Bucket",
			"DestinationBucketName": "my-logging-bucket",
		},
	})

	not deny["logging destination bucket is not correct: MyS3Bucket"] with input as inp
	not deny["logging prefix is not set correctly for bucket: MyS3Bucket"] with input as inp
}
