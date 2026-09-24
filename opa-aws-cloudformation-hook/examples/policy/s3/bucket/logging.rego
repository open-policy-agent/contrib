package aws.s3.bucket

import rego.v1

deny contains msg if {
	not bucket_excluded_logging
	not valid_logging_prefix
	msg := sprintf("logging prefix is not set correctly for bucket: %s", [input.resource.id])
}

deny contains msg if {
	not bucket_excluded_logging
	not valid_logging_destination
	msg := sprintf("logging destination bucket is not correct: %s", [input.resource.id])
}

valid_logging_prefix if {
	prefix := input.resource.properties.LoggingConfiguration.LogFilePrefix
	prefix == concat("-", ["s3-logs", input.resource.properties.BucketName])
}

valid_logging_prefix if {
	not input.resource.properties.BucketName
	input.resource.properties.LoggingConfiguration.LogFilePrefix == "s3-logs"
}

valid_logging_destination if {
	input.resource.properties.LoggingConfiguration.DestinationBucketName == "my-logging-bucket"
}

bucket_excluded_logging if {
	some prefix in {"excluded-", "access-", "secure-"}
	startswith(input.resource.properties.BucketName, prefix)
}
