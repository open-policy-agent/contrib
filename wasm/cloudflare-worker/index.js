const Rego = require("@open-policy-agent/opa-wasm")

// REGO_WASM is the resource the compiled policy is loaded into
var policy_wasm = REGO_WASM

addEventListener('fetch', event => {
    event.respondWith(handleRequest(event.request))
})

// Load WASM compiled policy, the loading is done asynchronously.
var rego = new Rego()
var loaded_policy = null
rego.load_policy(policy_wasm).then(policy => {
    loaded_policy = policy
}, error => {
    console.error("failed to load policy: " + error)
})

async function handleRequest(request) {
    //console.time("eval")
    
    // The policy may not have been loaded yet..
    // until then deny everything
    if (loaded_policy == null) {
        return new Response('{"error": "Policy not ready yet."}',
        { status: 503, statusText: "Service Unavailable" })
    }

    // the Request object doesn't have a "path"
    // field, only "url". So we add it ourselves
    url = new URL(request.url)
    request.path = url.pathname

    input_json = JSON.stringify(request)
    
    //console.log(input_json)

    allow = loaded_policy.eval_bool(input_json)
    
    //console.log("allow = " + allow)

    if (!allow) {
        // Short circuit the request here.
        return new Response('{"error": "Not allowed by policy"}',
        { status: 403, statusText: "Forbidden" })
    }
    
    // Allowed request
    
    const response = await fetch(request)

    //console.timeEnd("eval")
    return response
}
