package aws.rds.dbinstance

import rego.v1

deny contains msg if {
	not valid_storage_encryption
	msg := sprintf("storage encryption not enabled for: %s", [input.resource.id])
}

valid_storage_encryption if input.resource.properties.StorageEncrypted == "true"
