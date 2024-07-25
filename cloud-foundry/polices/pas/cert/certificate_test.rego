package cert_test

import rego.v1

import data.cert

test_expired_cert_is_expired if {
	expiry_date := ["2019-09-13T15:16:21Z"]
	actual := cert.determine_if_expired(expiry_date)
	expected := [{
		"date": "2019-09-13T15:16:21Z",
		"expired": true,
	}]

	actual == expected
}

test_non_expired_cert_is_not_expired if {
	expiry_date := ["2050-09-13T15:16:21Z"]
	actual := cert.determine_if_expired(expiry_date)
	expected := []

	actual == expected
}

test_get_certificate_expiry if {
	"2019-10-13T15:16:21Z" in cert.expiry(fake_cert)
}

test_get_certificate_expiry_two if {
	"2018-03-18T15:40:19Z" in cert.expiry(fake_cert2)
}

test_get_multiple_expiry if {
	mock_certs := [fake_cert, fake_cert2]
	expected := [
		"2019-10-13T15:16:21Z",
		"2018-03-18T15:40:19Z",
	]

	actual := cert.expiry(concat("", mock_certs))

	actual == expected
}

test_expiring_dates_return_true if {
	expiry_date := ["2019-09-13T15:16:21Z", "2019-10-13T15:16:21Z"]
	actual := cert.determine_if_expired(expiry_date)
	expected := [
		{
			"date": "2019-09-13T15:16:21Z",
			"expired": true,
		},
		{
			"date": "2019-10-13T15:16:21Z",
			"expired": true,
		},
	]

	actual == expected
}

test_cert_separation if {
	expected := [
		fake_cert,
		fake_cert2,
	]
	actual := cert.separate_certs(cert_chain)

	actual == expected
}

test_fail_on_no_certs if {
	actual := cert.expiry("meow")
	expected := 0

	count(actual) == expected
}

fake_cert := `-----BEGIN CERTIFICATE-----
MIIESjCCAzKgAwIBAgIJAJpBvIaxZbsqMA0GCSqGSIb3DQEBCwUAMIGYMQswCQYD
VQQGEwJVUzELMAkGA1UECAwCTlkxETAPBgNVBAcMCE5ldyBZb3JrMRQwEgYDVQQK
DAtNaXR0ZW5zIE9yZzEWMBQGA1UECwwNTGFicyBQbGF0Zm9ybTEZMBcGA1UEAwwQ
d2ViLmN1c3RvbWVyLm9yZzEgMB4GCSqGSIb3DQEJARYRZmFrZXlAbWNmYWtlcy5j
b20wHhcNMTkwOTEzMTUxNjIxWhcNMTkxMDEzMTUxNjIxWjCBmDELMAkGA1UEBhMC
VVMxCzAJBgNVBAgMAk5ZMREwDwYDVQQHDAhOZXcgWW9yazEUMBIGA1UECgwLTWl0
dGVucyBPcmcxFjAUBgNVBAsMDUxhYnMgUGxhdGZvcm0xGTAXBgNVBAMMEHdlYi5j
dXN0b21lci5vcmcxIDAeBgkqhkiG9w0BCQEWEWZha2V5QG1jZmFrZXMuY29tMIIB
IjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsr2TolHBt18D9ihtN+vvksLu
v+O8qpfzSbpj6WdOSNCsUSQiCVOeZMChoq4Lm7e7WHruoCj/el7+FxW6To/OkJwh
LSkfCBDYqqGq7AwTLuZ2/Hu/l93mqPkHGCQzYAu0R4NCh97uWa5PKdMqJKf7Bmrp
eMyLbsjEpImVw6hnEK7zllrzZDNLEXxhlakLw9TV2VfghMdo7TKqPYvnmzu/n+Zz
wKK0FQbc0YJoka/tMm/qe1GOnTSNIr7vp+ovPGh46/VU36YiKAiKfbWn3X6EPqSu
JabD55LBUQQRrrsLlKgoNOWbjyG73fJ8xd0/E8vtDyj9JgidfOikKcR33F/O/wID
AQABo4GUMIGRMA4GA1UdDwEB/wQEAwIFoDB0BgNVHREEbTBrghYqLnN5cy53ZWIu
Y3VzdG9tZXIub3JnghcqLmFwcHMud2ViLmN1c3RvbWVyLm9yZ4IcKi5sb2dpbi5z
eXMud2ViLmN1c3RvbWVyLm9yZ4IaKi51YWEuc3lzLndlYi5jdXN0b21lci5vcmcw
CQYDVR0TBAIwADANBgkqhkiG9w0BAQsFAAOCAQEAhUuD26d2hCfrq3OytwZMTa23
uZVKHxM51EnovPwUHDIAcNYAF3pkcn4eEZI4hT0aQMq5O27r6wVwSTWywHZLw378
l5UwYBXUXtbuAyGXDgKDAJGUcGQ1Z7fQ0YtXXMUAutxNpc54hjaEdgAfyxk3kOmT
/k6mWGzaYkxGmuxwh4uM831fEeQaKjBdlBuPggUa0ZKpF/6J/vFGCiqz7cfI1Pde
xuaRoq5on8aUDrHW5wjGwxLncivpyDP2Pt3FEKw5tGlQCU0+b3JZBEYmAfgENkEf
K1Qp01rVVJ0+0lPI7M5cDEQhJtEg2PXzzy/inGK28JcLmNXY88BP5WLPou1JGw==
-----END CERTIFICATE-----
`

fake_cert2 := `-----BEGIN CERTIFICATE-----
MIICMzCCAZygAwIBAgIJALiPnVsvq8dsMA0GCSqGSIb3DQEBBQUAMFMxCzAJBgNV
BAYTAlVTMQwwCgYDVQQIEwNmb28xDDAKBgNVBAcTA2ZvbzEMMAoGA1UEChMDZm9v
MQwwCgYDVQQLEwNmb28xDDAKBgNVBAMTA2ZvbzAeFw0xMzAzMTkxNTQwMTlaFw0x
ODAzMTgxNTQwMTlaMFMxCzAJBgNVBAYTAlVTMQwwCgYDVQQIEwNmb28xDDAKBgNV
BAcTA2ZvbzEMMAoGA1UEChMDZm9vMQwwCgYDVQQLEwNmb28xDDAKBgNVBAMTA2Zv
bzCBnzANBgkqhkiG9w0BAQEFAAOBjQAwgYkCgYEAzdGfxi9CNbMf1UUcvDQh7MYB
OveIHyc0E0KIbhjK5FkCBU4CiZrbfHagaW7ZEcN0tt3EvpbOMxxc/ZQU2WN/s/wP
xph0pSfsfFsTKM4RhTWD2v4fgk+xZiKd1p0+L4hTtpwnEw0uXRVd0ki6muwV5y/P
+5FHUeldq+pgTcgzuK8CAwEAAaMPMA0wCwYDVR0PBAQDAgLkMA0GCSqGSIb3DQEB
BQUAA4GBAJiDAAtY0mQQeuxWdzLRzXmjvdSuL9GoyT3BF/jSnpxz5/58dba8pWen
v3pj4P3w5DoOso0rzkZy2jEsEitlVM2mLSbQpMM+MUVQCQoiG6W9xuCFuxSrwPIS
pAqEAuV4DNoxQKKWmhVv+J0ptMWD25Pnpxeq5sXzghfJnslJlQND
-----END CERTIFICATE-----
`

cert_chain := concat("", [fake_cert, fake_cert2])
