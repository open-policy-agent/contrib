package credhub

import rego.v1

deny_if_not_exactly_one_primary contains msg if {
	keys := [key.primary |
		some key in input["product-properties"][".properties.credhub_key_encryption_passwords"].value
	]

	count(keys) != 1

	msg := sprintf("Must have exactly one primary encryption key for credhub, found %d", [count(keys)])
}

deny_not_enough_chars contains msg if {
	some val in input["product-properties"][".properties.credhub_key_encryption_passwords"].value
	val.primary

	count(val.key.secret) < 20

	msg := sprintf("Primary key must be at least 20 characters, found %v", [count(val.key.secret)])
}
