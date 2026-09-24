package aws.rds.dbinstance_test

import rego.v1

import data.aws.rds.dbinstance.deny

import data.test_helpers.create_with_properties

test_deny_storage_encryption_disabled if {
	inp := create_with_properties("AWS::RDS::DBInstance", "RDSInstance", {"StorageEncrypted": "false"})

	deny["storage encryption not enabled for: RDSInstance"] with input as inp
}

test_deny_storage_encryption_not_set if {
	inp := create_with_properties("AWS::RDS::DBInstance", "RDSInstance", {"properties": {
		"DBName": "OPA-DB",
		"Engine": "MySQL",
	}})

	deny["storage encryption not enabled for: RDSInstance"] with input as inp
}

test_allow_storage_encryption_enabled if {
	inp := create_with_properties("AWS::RDS::DBInstance", "RDSInstance", {"StorageEncrypted": "true"})

	count(deny) == 0 with input as inp
}
