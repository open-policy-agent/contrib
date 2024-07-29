package main

import rego.v1

find(json, desired_value) := path if {
	some path
	walk(json, [path, desired_value])
}

deny_if_ciphers_missing contains msg if {
	desired_cipher := "ECDHE-RSA-AES128-GCM-SHA256:TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"
	route_to_value := find(input, desired_cipher)

	# count(routeToValue) < 1

	msg := sprintf(
		"expected cipher configuration of: %v\n please update the value following this json path: %v",
		[desired_cipher, route_to_value],
	)
}

# the above is closer to what i want
# now I'm getting the path to the output
