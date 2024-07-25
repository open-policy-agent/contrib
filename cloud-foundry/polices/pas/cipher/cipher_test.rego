package main_test

import rego.v1

import data.main

test_if_ciphers_match if {
	val := "ECDHE-RSA-AES128-GCM-SHA256:TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"
	obj := {".properties.gorouter_ssl_ciphers": {"value": val}}

	result := main.deny_if_ciphers_missing with input as obj

	expect := concat("\n", [
		"expected cipher configuration of: ECDHE-RSA-AES128-GCM-SHA256:TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
		" please update the value following this json path: [\".properties.gorouter_ssl_ciphers\", \"value\"]",
	])

	result[expect]
}
