package credhub

deny_if_not_exactly_one_primary[msg] {
        keys := [key |
                key := input["product-properties"][".properties.credhub_key_encryption_passwords"].value[i].primary
        ]

        keys != [true]
        msg = sprintf("Must have exactly one primary encryption key for credhub, found %d", [count(keys)])
}

deny_not_enough_chars[msg] {
    some i
    input["product-properties"][".properties.credhub_key_encryption_passwords"].value[i].primary
    key := input["product-properties"][".properties.credhub_key_encryption_passwords"].value[i].key.secret
    count(key) < 20
    msg := sprintf("Primary key must be at least 20 characters, found %v", [count(key)])
}