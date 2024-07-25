package credhub_test

import rego.v1

import data.credhub

test_deny_credhub_suite if {
	credhub.deny_if_not_exactly_one_primary with input as pas_key
	credhub.deny_not_enough_chars with input as pas_key
}

pas_key := {"product-properties": {
	".properties.credhub_hsm_provider_partition_password": {"value": [{"primary": false}]},
	".properties.credhub_key_encryption_passwords": {"value": [{"primary": "1234567890123456789"}]},
}}
