# This file defines the entry point for the image enforcer policy.
#
# The default policy will not reject any images unless they have been blacklisted.
# By default, the blacklist is empty.

package io.k8s.image_policy

import rego.v1

verify := {
	"apiVersion": "imagepolicy.k8s.io/v1alpha1",
	"kind": "ImageReview",
	"status": {"allowed": allow},
}

default allow := false

allow if not deny

default deny := false
