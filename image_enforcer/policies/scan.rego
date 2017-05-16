# Implements blacklisting of known-vulnerable images.

package io.k8s.image_policy

secure_namespaces = {
    "production",
}

deny {
    secure_namespaces[input.spec.namespace]
    image_layers[tuple]
    tuple = [org, name, tag, layer]
    vulnerable_layers[layer]
}

image_layers[[org, name, tag, layer]] {
    input.spec.containers[i].image = img
    split(img, ":", [repo, tag])
    split(repo, "/", [org, name])
    data.docker.layers[org][name][tag].layers[_].digest = layer
}

vulnerable_layers[layer] {
    data.clair.layers[layer].Layer.Features[_].Vulnerabilities != []
}
