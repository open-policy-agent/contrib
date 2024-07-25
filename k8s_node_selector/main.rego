package system

import rego.v1

# The top level response sent to the webrequest
main := {
	"apiVersion": "admission.k8s.io/v1beta1",
	"kind": "AdmissionReview",
	"response": response,
}

# If the conditions on the `response` below aren't met this default `allow` response
#  is returned.
default response := {"allowed": true}

# These are the options for the response body sent back to admissions controller request
#   it starts with a number of conditions which have to be met for it to take effect
#
# regal ignore:rule-length
response := output if {
	#
	# Option 1: Return Allowed=false if
	#  - A Pod
	#  - In a namespace we care about
	#  - With nodeselector already specified <- Shouldn't be done for a namespace with a mapping
	is_pod
	pool_for_namespace(input.request.namespace)
	has_node_selector

	output := {
		"allowed": false,
		"status": {
			"code": 403,
			"reason": "Manually specifying NodeSelector not supported in this namespace.",
		},
	}
} else := output if {
	#
	# Option 2: Set the AgentPool=CorrectPool
	#  - A Pod
	#  - Without a nodeselector already specified
	#  - In a namespace we care about
	#  - Has a mapping set for the namespace
	is_pod
	not has_node_selector

	# Generate the JSON Patch object
	patch := {
		"op": "add",
		"path": "/spec/nodeSelector",
		"value": {"agentpool": pool_for_namespace(input.request.namespace)},
	}

	# Patches have to be an array of base64 encoded JSON Patches so lets
	#  make our single patch into an array, serialize as JSON and base64 encode.
	patches := [patch]
	patch_encoded := base64.encode(json.marshal(patches))

	# Output a trace use `opa test *.rego -v --explain full` to see them.
	trace(sprintf("POLICY:generatedPatch raw = '%s'", [patches])) # regal ignore:print-or-trace-call
	trace(sprintf("POLICY:generatedPatch encoded = '%s'", [patch_encoded])) # regal ignore:print-or-trace-call

	# Generate the patch response and return it! We're done!
	output := {
		"allowed": true,
		"patchType": "JSONPatch",
		"patch": patch_encoded,
	}
}

# Rule: Check if the item submitted is a pod.
is_pod if {
	input.request.kind.kind == "Pod"
}

# Rule: Check if pod already has a `nodeSelector` set
has_node_selector if {
	input.request.object.spec.nodeSelector
	count(input.request.object.spec.nodeSelector) > 0
}

# Rule: For the given namespace get the `agentpool` label set
# on the namespace object itself

pool_for_namespace(namespace) := pool_label if {
	# regal ignore:print-or-trace-call,external-reference
	trace(sprintf("POLICY:namespace raw = '%s'", [data.kubernetes.namespaces]))

	# regal ignore:external-reference
	cluster_namespace := data.kubernetes.namespaces[namespace]

	pool_label = cluster_namespace.metadata.labels.agentpool
}
