package cert

parse_certificate(cert) = parsedCertificate {
        strippedCert := replace(replace(cert, "-----END CERTIFICATE-----", ""), "-----BEGIN CERTIFICATE-----", "")
        parsedCertificate := crypto.x509.parse_certificates(strippedCert)
}

separate_certs(certChain) = cleanedCerts {
        addDelimeter := replace(certChain, "-----END CERTIFICATE-----\n", "-----END CERTIFICATE-----\n&&&&")
        splitCerts := split(addDelimeter, "&&&&")
        count(splitCerts) > 0

        cleanedCerts := array.slice(splitCerts, 0, count(splitCerts) - 1)
}

get_certificate_expiry(rawCertChain) = expiryDate {
        certArray := separate_certs(rawCertChain)
        parsedCerts := [parsedCert |
                cert := certArray[_]
                parsedCert := parse_certificate(cert)
        ]

        expiryDate := [expiryDate |
                expiryDate := parsedCerts[_][_].NotAfter
        ]
}

determine_if_expired(dates) = certsForRenewal {
        thirty_days_in_nanoseconds := 2.592e+15

        currentTime_nano := time.now_ns()
        certsForRenewal := [expired |
                date := dates[_]
                certExpiryDate_nano := time.parse_rfc3339_ns(date)
                timeDelta := certExpiryDate_nano - currentTime_nano
                timeDelta <= thirty_days_in_nanoseconds

                expired := {
                        "date": date,
                        "expired": timeDelta <= thirty_days_in_nanoseconds,
                }
        ]
}

deny_certs_not_present[msg] {
        exists := [certs |
                certs := input.certs
        ] #you will need to provide a path to a cert

        count(exists) == 0

        msg = sprintf("No certs in provided, either in path or input object: %v", [exists])
}

deny_thirty_days[msg] {
        # must manually define path to cert. JSON input
        # key values are accessed using bracket notation rather than dot "." notation
        certs := input.certs #you will need to provide a path to a cert
        expirys := get_certificate_expiry(certs)
        isExpired := determine_if_expired(expirys)

        count(isExpired) > 0

        msg = sprintf("Your certificate expires on this date %v please update cert", [isExpired])
}