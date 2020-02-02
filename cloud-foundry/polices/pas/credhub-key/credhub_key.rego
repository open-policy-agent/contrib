package credhub

deny_if_not_exactly_one_primary[msg] {
        keys := [key |
                input["product-properties"][".properties.credhub_key_encryption_passwords"].value[i].primary # only do the next statement if this is true
                key := input["product-properties"][".properties.credhub_key_encryption_passwords"].value[i].primary
        ]

        keys != [true]
        msg = sprintf("Must have exactly one primary encryption key for credhub, found %d", [count(keys)])
}

deny_not_enough_chars[msg] {
        key := [key |
                input["product-properties"][".properties.credhub_key_encryption_passwords"].value[i].primary # only do the next statement if this is true
                key := input["product-properties"][".properties.credhub_key_encryption_passwords"].value[i].key.secret
        ]

        count(key[0]) < 20

        msg = sprintf("Primary key must be at least 20 characters, found %v", [count(key[0])])
}