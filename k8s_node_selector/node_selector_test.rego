package system_test

import rego.v1

import data.node_selector.testdata

import data.system

# -----------------------------------------------------------
# Test: Test patch applied to pods
# -----------------------------------------------------------

# patch we expect to be generated
expected_patch := {
	"op": "add",
	"path": "/spec/nodeSelector",
	"value": {"agentpool": "pool1"},
}

# Checks the response is a patch response
is_patch_response(res) if {
	# Is the response patch type correct?
	res.response.patchType == "JSONPatch"

	# Is the patch body set?
	res.response.patch

	# Is the patch body an array with at least one item?
	count(res.response.patch) > 0
}

# Test that the controller correctly sets a patch on a pod
# to assign it the correct `nodeSelector` of `agentpool=pool1`
# when in a namespace with label `agentpool=pool1`
test_response_patch_valid_mapping if {
	# Invoke the policy main with a pod which doesn't have a node selector
	#  and is in the default namespace
	body := system.main with input as testdata.example_pod_doesnt_have_node_selector
		with data.kubernetes.namespaces as testdata.example_namespaces_with_label

	# Check policy returned an allowed response
	body.response.allowed == true

	# Check the response is a patch response
	is_patch_response(body)

	# The admission controller response is an array of base64 encoded
	# jsonpatches so deserialize so we can review them.
	patches := json.unmarshal(base64.decode(body.response.patch))

	# Output some tracing... `opa test *.rego -v --explain full` to see them
	trace(sprintf("TEST:appliedPatch = '%s'", [patches])) # regal ignore:print-or-trace-call
	trace(sprintf("TEST:expectedPatch = '%s'", [expected_patch])) # regal ignore:print-or-trace-call

	# Check the policy created the expected patch
	expected_patch in patches
}

# Test that the controller accepts a pod in a namespace without an agent pool mapping
test_response_patch_invalid_mapping if {
	# Invoke the policy main with a pod which doesn't have a node selector
	#  and is in the default namespace
	body := system.main with input as testdata.example_pod_has_nonexisting_namespace
		with data.kubernetes.namespaces as testdata.example_namespaces_with_label

	# Check policy returned an allowed response
	body.response.allowed == true
}

# Test that the controller rejects a pod with a nodeselector already set
test_response_patch_valid_mapping_with_nodeselector if {
	# Invoke the policy main with a pod which doesn't have a node selector
	#  and is in the default namespace
	body := system.main with input as testdata.example_pod_has_node_selector
		with data.kubernetes.namespaces as testdata.example_namespaces_with_label

	# Check policy returned an allowed response
	body.response.allowed == false
}

# -----------------------------------------------------------
# Helpers: Test helper functions used by the policy
# -----------------------------------------------------------

test_has_selector_has_selector if {
	system.has_node_selector with input as testdata.example_pod_has_node_selector
}

test_has_selector_doesnt_have_selector if {
	not system.has_node_selector with input as testdata.example_pod_doesnt_have_node_selector
}

test_response_allowed if {
	body := system.main with input as testdata.example_pod_doesnt_have_node_selector
	body.response.allowed == true
}

test_is_pod_true if {
	system.is_pod with input as testdata.example_pod_has_node_selector
}
