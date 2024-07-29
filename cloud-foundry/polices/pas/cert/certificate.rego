package cert

import rego.v1

parse_certificate(cert) := parsed_certificate if {
	stripped_cert := replace(replace(cert, "-----END CERTIFICATE-----", ""), "-----BEGIN CERTIFICATE-----", "")
	parsed_certificate := crypto.x509.parse_certificates(stripped_cert)
}

separate_certs(cert_chain) := cleaned_certs if {
	add_delimeter := replace(cert_chain, "-----END CERTIFICATE-----\n", "-----END CERTIFICATE-----\n&&&&")
	split_certs := split(add_delimeter, "&&&&")
	count(split_certs) > 0

	cleaned_certs := array.slice(split_certs, 0, count(split_certs) - 1)
}

expiry(raw_cert_chain) := expiry_date if {
	cert_array := separate_certs(raw_cert_chain)
	parsed_certs := [parsed_cert |
		some cert in cert_array
		parsed_cert := parse_certificate(cert)
	]

	expiry_date := [expiry_date |
		expiry_date := parsed_certs[_][_].NotAfter
	]
}

determine_if_expired(dates) := certs_for_renewal if {
	thirty_days_in_nanoseconds := 2.592e+15

	certs_for_renewal := [expired |
		some date in dates
		cert_expiry_date_nano := time.parse_rfc3339_ns(date)
		time_delta := cert_expiry_date_nano - time.now_ns()
		time_delta <= thirty_days_in_nanoseconds

		expired := {
			"date": date,
			"expired": time_delta <= thirty_days_in_nanoseconds,
		}
	]
}

deny_certs_not_present contains msg if {
	exists := [certs |
		certs := input.certs
	] # you will need to provide a path to a cert

	count(exists) == 0

	msg = sprintf("No certs in provided, either in path or input object: %v", [exists])
}

deny_thirty_days contains msg if {
	# must manually define path to cert. JSON input
	# key values are accessed using bracket notation rather than dot "." notation
	certs := input.certs # you will need to provide a path to a cert
	expirys := expiry(certs)
	is_expired := determine_if_expired(expirys)

	count(is_expired) > 0

	msg = sprintf("Your certificate expires on this date %v please update cert", [is_expired])
}
