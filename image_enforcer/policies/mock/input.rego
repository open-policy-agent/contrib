# This file contains test inputs for the image enforcer policies.

package mock

import rego.v1

allowed := {
	"kind": "ImageReview",
	"apiVersion": "imagepolicy.k8s.io/v1alpha1",
	"spec": {
		"containers": [{"image": "openpolicyagent/opa:0.4.9"}],
		"namespace": "production",
	},
}

denied := {
	"kind": "ImageReview",
	"apiVersion": "imagepolicy.k8s.io/v1alpha1",
	"spec": {
		"containers": [
			{"image": "openpolicyagent/opa:0.4.9"},
			{"image": "library/wordpress:4-apache"},
		],
		"namespace": "production",
	},
}
