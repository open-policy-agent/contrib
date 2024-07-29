# Implements blacklisting of known-vulnerable images.

package io.k8s.image_policy

import rego.v1

secure_namespaces := {"production"}

deny if {
	input.spec.namespace in secure_namespaces
	some tuple in image_layers
	some [org, name, tag, layer] in tuple
	layer in vulnerable_layers
}

image_layers contains [org, name, tag, layer.digest] if {
	some container in input.spec.containers
	[repo, tag] := split(container.image, ":")
	[org, name] := split(repo, "/")
	some layer in data.docker.layers[org][name][tag].layers
}

vulnerable_layers contains layer if {
	some layer
	data.clair.layers[layer].Layer.Features[_].Vulnerabilities != [] # regal ignore:not-equals-in-loop
}
