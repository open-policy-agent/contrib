package main

find(json, desiredValue) = route {
    some path
    walk(json, [path, desiredValue])


    count(path) != 0
    route := path
}

deny_if_ciphers_missing[msg] {
    desiredCipher := "ECDHE-RSA-AES128-GCM-SHA256:TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"
    routeToValue := find(input, desiredCipher)

    # count(routeToValue) < 1
    true

    msg = sprintf("expected cipher configuration of: %v\n please update the value following this json path: %v", [desiredCipher, routeToValue])
}

#the above is closer to what i want
#now I'm getting the path to the output