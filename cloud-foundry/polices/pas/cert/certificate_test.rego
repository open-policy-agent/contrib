package cert

test_expired_cert_is_expired {

    expiryDate := ["2019-09-13T15:16:21Z"]
    actual := determine_if_expired(expiryDate)
    expected := [{
        "date": "2019-09-13T15:16:21Z",
        "expired": true
    }]

    actual == expected
}

test_non_expired_cert_is_not_expired {

    expiryDate := ["2050-09-13T15:16:21Z"]
    actual := determine_if_expired(expiryDate)
    expected := []

    actual == expected
}

test_get_certificate_expiry {
    actual := get_certificate_expiry(fakeCert)
    expected := "2019-10-13T15:16:21Z"

    actual[_] == expected
    
}

test_get_certificate_expiry_two {
    actual := get_certificate_expiry(fakeCert2)
    expected := "2018-03-18T15:40:19Z"

    actual[_] == expected
}

test_get_multiple_expiry {
    mockCerts := [fakeCert, fakeCert2]
    expected := [
            "2019-10-13T15:16:21Z", 
            "2018-03-18T15:40:19Z",
    ]

    actual := get_certificate_expiry(concat("", mockCerts))

    actual == expected
}

test_expiring_dates_return_true {

    expiryDate := ["2019-09-13T15:16:21Z", "2019-10-13T15:16:21Z"]
    actual := determine_if_expired(expiryDate)
    expected := [
        {
            "date": "2019-09-13T15:16:21Z",
            "expired": true
        },
        {
            "date": "2019-10-13T15:16:21Z",
            "expired": true
        }
    ]

    actual == expected
}

test_cert_separation {
    mockCerts := certChain
    expected := [
            fakeCert, 
            fakeCert2,
    ]
    actual := separate_certs(mockCerts)

    actual == expected

}

test_fail_on_no_certs {
    actual := get_certificate_expiry("meow")
    expected := 0

    count(actual) == expected
}

fakeCert ="-----BEGIN CERTIFICATE-----\nMIIESjCCAzKgAwIBAgIJAJpBvIaxZbsqMA0GCSqGSIb3DQEBCwUAMIGYMQswCQYD\nVQQGEwJVUzELMAkGA1UECAwCTlkxETAPBgNVBAcMCE5ldyBZb3JrMRQwEgYDVQQK\nDAtNaXR0ZW5zIE9yZzEWMBQGA1UECwwNTGFicyBQbGF0Zm9ybTEZMBcGA1UEAwwQ\nd2ViLmN1c3RvbWVyLm9yZzEgMB4GCSqGSIb3DQEJARYRZmFrZXlAbWNmYWtlcy5j\nb20wHhcNMTkwOTEzMTUxNjIxWhcNMTkxMDEzMTUxNjIxWjCBmDELMAkGA1UEBhMC\nVVMxCzAJBgNVBAgMAk5ZMREwDwYDVQQHDAhOZXcgWW9yazEUMBIGA1UECgwLTWl0\ndGVucyBPcmcxFjAUBgNVBAsMDUxhYnMgUGxhdGZvcm0xGTAXBgNVBAMMEHdlYi5j\ndXN0b21lci5vcmcxIDAeBgkqhkiG9w0BCQEWEWZha2V5QG1jZmFrZXMuY29tMIIB\nIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsr2TolHBt18D9ihtN+vvksLu\nv+O8qpfzSbpj6WdOSNCsUSQiCVOeZMChoq4Lm7e7WHruoCj/el7+FxW6To/OkJwh\nLSkfCBDYqqGq7AwTLuZ2/Hu/l93mqPkHGCQzYAu0R4NCh97uWa5PKdMqJKf7Bmrp\neMyLbsjEpImVw6hnEK7zllrzZDNLEXxhlakLw9TV2VfghMdo7TKqPYvnmzu/n+Zz\nwKK0FQbc0YJoka/tMm/qe1GOnTSNIr7vp+ovPGh46/VU36YiKAiKfbWn3X6EPqSu\nJabD55LBUQQRrrsLlKgoNOWbjyG73fJ8xd0/E8vtDyj9JgidfOikKcR33F/O/wID\nAQABo4GUMIGRMA4GA1UdDwEB/wQEAwIFoDB0BgNVHREEbTBrghYqLnN5cy53ZWIu\nY3VzdG9tZXIub3JnghcqLmFwcHMud2ViLmN1c3RvbWVyLm9yZ4IcKi5sb2dpbi5z\neXMud2ViLmN1c3RvbWVyLm9yZ4IaKi51YWEuc3lzLndlYi5jdXN0b21lci5vcmcw\nCQYDVR0TBAIwADANBgkqhkiG9w0BAQsFAAOCAQEAhUuD26d2hCfrq3OytwZMTa23\nuZVKHxM51EnovPwUHDIAcNYAF3pkcn4eEZI4hT0aQMq5O27r6wVwSTWywHZLw378\nl5UwYBXUXtbuAyGXDgKDAJGUcGQ1Z7fQ0YtXXMUAutxNpc54hjaEdgAfyxk3kOmT\n/k6mWGzaYkxGmuxwh4uM831fEeQaKjBdlBuPggUa0ZKpF/6J/vFGCiqz7cfI1Pde\nxuaRoq5on8aUDrHW5wjGwxLncivpyDP2Pt3FEKw5tGlQCU0+b3JZBEYmAfgENkEf\nK1Qp01rVVJ0+0lPI7M5cDEQhJtEg2PXzzy/inGK28JcLmNXY88BP5WLPou1JGw==\n-----END CERTIFICATE-----\n"
fakeCert2 = "-----BEGIN CERTIFICATE-----\nMIICMzCCAZygAwIBAgIJALiPnVsvq8dsMA0GCSqGSIb3DQEBBQUAMFMxCzAJBgNV\nBAYTAlVTMQwwCgYDVQQIEwNmb28xDDAKBgNVBAcTA2ZvbzEMMAoGA1UEChMDZm9v\nMQwwCgYDVQQLEwNmb28xDDAKBgNVBAMTA2ZvbzAeFw0xMzAzMTkxNTQwMTlaFw0x\nODAzMTgxNTQwMTlaMFMxCzAJBgNVBAYTAlVTMQwwCgYDVQQIEwNmb28xDDAKBgNV\nBAcTA2ZvbzEMMAoGA1UEChMDZm9vMQwwCgYDVQQLEwNmb28xDDAKBgNVBAMTA2Zv\nbzCBnzANBgkqhkiG9w0BAQEFAAOBjQAwgYkCgYEAzdGfxi9CNbMf1UUcvDQh7MYB\nOveIHyc0E0KIbhjK5FkCBU4CiZrbfHagaW7ZEcN0tt3EvpbOMxxc/ZQU2WN/s/wP\nxph0pSfsfFsTKM4RhTWD2v4fgk+xZiKd1p0+L4hTtpwnEw0uXRVd0ki6muwV5y/P\n+5FHUeldq+pgTcgzuK8CAwEAAaMPMA0wCwYDVR0PBAQDAgLkMA0GCSqGSIb3DQEB\nBQUAA4GBAJiDAAtY0mQQeuxWdzLRzXmjvdSuL9GoyT3BF/jSnpxz5/58dba8pWen\nv3pj4P3w5DoOso0rzkZy2jEsEitlVM2mLSbQpMM+MUVQCQoiG6W9xuCFuxSrwPIS\npAqEAuV4DNoxQKKWmhVv+J0ptMWD25Pnpxeq5sXzghfJnslJlQND\n-----END CERTIFICATE-----\n"
certChain = "-----BEGIN CERTIFICATE-----\nMIIESjCCAzKgAwIBAgIJAJpBvIaxZbsqMA0GCSqGSIb3DQEBCwUAMIGYMQswCQYD\nVQQGEwJVUzELMAkGA1UECAwCTlkxETAPBgNVBAcMCE5ldyBZb3JrMRQwEgYDVQQK\nDAtNaXR0ZW5zIE9yZzEWMBQGA1UECwwNTGFicyBQbGF0Zm9ybTEZMBcGA1UEAwwQ\nd2ViLmN1c3RvbWVyLm9yZzEgMB4GCSqGSIb3DQEJARYRZmFrZXlAbWNmYWtlcy5j\nb20wHhcNMTkwOTEzMTUxNjIxWhcNMTkxMDEzMTUxNjIxWjCBmDELMAkGA1UEBhMC\nVVMxCzAJBgNVBAgMAk5ZMREwDwYDVQQHDAhOZXcgWW9yazEUMBIGA1UECgwLTWl0\ndGVucyBPcmcxFjAUBgNVBAsMDUxhYnMgUGxhdGZvcm0xGTAXBgNVBAMMEHdlYi5j\ndXN0b21lci5vcmcxIDAeBgkqhkiG9w0BCQEWEWZha2V5QG1jZmFrZXMuY29tMIIB\nIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsr2TolHBt18D9ihtN+vvksLu\nv+O8qpfzSbpj6WdOSNCsUSQiCVOeZMChoq4Lm7e7WHruoCj/el7+FxW6To/OkJwh\nLSkfCBDYqqGq7AwTLuZ2/Hu/l93mqPkHGCQzYAu0R4NCh97uWa5PKdMqJKf7Bmrp\neMyLbsjEpImVw6hnEK7zllrzZDNLEXxhlakLw9TV2VfghMdo7TKqPYvnmzu/n+Zz\nwKK0FQbc0YJoka/tMm/qe1GOnTSNIr7vp+ovPGh46/VU36YiKAiKfbWn3X6EPqSu\nJabD55LBUQQRrrsLlKgoNOWbjyG73fJ8xd0/E8vtDyj9JgidfOikKcR33F/O/wID\nAQABo4GUMIGRMA4GA1UdDwEB/wQEAwIFoDB0BgNVHREEbTBrghYqLnN5cy53ZWIu\nY3VzdG9tZXIub3JnghcqLmFwcHMud2ViLmN1c3RvbWVyLm9yZ4IcKi5sb2dpbi5z\neXMud2ViLmN1c3RvbWVyLm9yZ4IaKi51YWEuc3lzLndlYi5jdXN0b21lci5vcmcw\nCQYDVR0TBAIwADANBgkqhkiG9w0BAQsFAAOCAQEAhUuD26d2hCfrq3OytwZMTa23\nuZVKHxM51EnovPwUHDIAcNYAF3pkcn4eEZI4hT0aQMq5O27r6wVwSTWywHZLw378\nl5UwYBXUXtbuAyGXDgKDAJGUcGQ1Z7fQ0YtXXMUAutxNpc54hjaEdgAfyxk3kOmT\n/k6mWGzaYkxGmuxwh4uM831fEeQaKjBdlBuPggUa0ZKpF/6J/vFGCiqz7cfI1Pde\nxuaRoq5on8aUDrHW5wjGwxLncivpyDP2Pt3FEKw5tGlQCU0+b3JZBEYmAfgENkEf\nK1Qp01rVVJ0+0lPI7M5cDEQhJtEg2PXzzy/inGK28JcLmNXY88BP5WLPou1JGw==\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\nMIICMzCCAZygAwIBAgIJALiPnVsvq8dsMA0GCSqGSIb3DQEBBQUAMFMxCzAJBgNV\nBAYTAlVTMQwwCgYDVQQIEwNmb28xDDAKBgNVBAcTA2ZvbzEMMAoGA1UEChMDZm9v\nMQwwCgYDVQQLEwNmb28xDDAKBgNVBAMTA2ZvbzAeFw0xMzAzMTkxNTQwMTlaFw0x\nODAzMTgxNTQwMTlaMFMxCzAJBgNVBAYTAlVTMQwwCgYDVQQIEwNmb28xDDAKBgNV\nBAcTA2ZvbzEMMAoGA1UEChMDZm9vMQwwCgYDVQQLEwNmb28xDDAKBgNVBAMTA2Zv\nbzCBnzANBgkqhkiG9w0BAQEFAAOBjQAwgYkCgYEAzdGfxi9CNbMf1UUcvDQh7MYB\nOveIHyc0E0KIbhjK5FkCBU4CiZrbfHagaW7ZEcN0tt3EvpbOMxxc/ZQU2WN/s/wP\nxph0pSfsfFsTKM4RhTWD2v4fgk+xZiKd1p0+L4hTtpwnEw0uXRVd0ki6muwV5y/P\n+5FHUeldq+pgTcgzuK8CAwEAAaMPMA0wCwYDVR0PBAQDAgLkMA0GCSqGSIb3DQEB\nBQUAA4GBAJiDAAtY0mQQeuxWdzLRzXmjvdSuL9GoyT3BF/jSnpxz5/58dba8pWen\nv3pj4P3w5DoOso0rzkZy2jEsEitlVM2mLSbQpMM+MUVQCQoiG6W9xuCFuxSrwPIS\npAqEAuV4DNoxQKKWmhVv+J0ptMWD25Pnpxeq5sXzghfJnslJlQND\n-----END CERTIFICATE-----\n"