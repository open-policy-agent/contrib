package system

import data.node_selector.testdata as testdata

#-----------------------------------------------------------
# Test: Test patch applied to pods
#-----------------------------------------------------------

# patch we expect to be generated
expectedPatch = {
	"op": "add",
    "path": "/spec/nodeSelector",
    "value": {
        "agentpool": "pool1",
    }
}

# Helper to check patch is set
hasPatch(patches, expectedPatch) {
    # One of the patches returned should match the `expectedPatch`
    patches[i] == expectedPatch
}

# Checks the response is a patch response
isPatchResponse(res) {
	# Is the response patch type correct?
    res.response.patchType == "JSONPatch"
    # Is the patch body set?
	res.response.patch
    # Is the patch body an array with at least one item?
    count(res.response.patch) > 0
}

# Test that the controller correctly sets a patch on a pod
#  to assign it the correct `nodeSelector` of `agentpool=pool1`
#  when in a namespace with label `agentpool=pool1`
test_response_patch_valid_mapping {
    # Invoke the policy main with a pod which doesn't have a node selector
    #  and is in the default namespace
    body := main 
        with input as testdata.example_pod_doesnt_have_node_selector 
        with data.kubernetes.namespaces as testdata.example_namespaces_with_label
    
    # Check policy returned an allowed response
    body.response.allowed == true

    # Check the response is a patch response
    isPatchResponse(body)

    # The admission controller response is an array of base64 encoded
    # jsonpatches so deserialize so we can review them. 
    patches := json.unmarshal(base64.decode(body.response.patch))

    # Output some tracing... `opa test *.rego -v --explain full` to see them
    trace(sprintf("TEST:appliedPatch = '%s'", [patches]))
    trace(sprintf("TEST:expectedPatch = '%s'", [expectedPatch]))

    # Check the policy created the expected patch
    hasPatch(patches, expectedPatch)
}


# Test that the controller accepts a pod in a namespace without an agent pool mapping
test_response_patch_invalid_mapping {
    # Invoke the policy main with a pod which doesn't have a node selector
    #  and is in the default namespace
    body := main 
        with input as testdata.example_pod_has_nonexisting_namespace
        with data.kubernetes.namespaces as testdata.example_namespaces_with_label
    
    # Check policy returned an allowed response
    body.response.allowed == true
}

# Test that the controller rejects a pod with a nodeselector already set
test_response_patch_valid_mapping_with_nodeselector {
    # Invoke the policy main with a pod which doesn't have a node selector
    #  and is in the default namespace
    body := main 
        with input as testdata.example_pod_has_node_selector 
        with data.kubernetes.namespaces as testdata.example_namespaces_with_label
    
    # Check policy returned an allowed response
    body.response.allowed == false
}

#-----------------------------------------------------------
# Helpers: Test helper functions used by the policy
#-----------------------------------------------------------

test_hasSelector_has_selector {
    hasNodeSelector with input as testdata.example_pod_has_node_selector
}

test_hasSelector_doesnt_have_selector {
    not hasNodeSelector with input as testdata.example_pod_doesnt_have_node_selector
}

test_response_allowed {
    body := main with input as testdata.example_pod_doesnt_have_node_selector
    body.response.allowed == true
}

test_isPod_true {
    isPod with input as testdata.example_pod_has_node_selector
}