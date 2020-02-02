package main

test_if_ciphers_match {
    deny_if_ciphers_missing[_] with input as { 
            ".properties.gorouter_ssl_ciphers": {
                    "value": "ECDHE-RSA-AES128-GCM-SHA256:TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
                }
            }
}