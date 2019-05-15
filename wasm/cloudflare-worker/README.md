# WASM Policy Worker
This example shows making a [Cloudflare Worker](https://www.cloudflare.com/products/cloudflare-workers/)
which can enforce Rego Policies that have been compiled into WebAssembly (wasm).

## The Worker

The worker is a single javascript file that uses some `npm` modules, most importantly
the [@open-policy-agent/opa-wasm](https://github.com/open-policy-agent/npm-opa-wasm)
module.

Because the Cloudflare worker should be a single self-contained script we will use
[Webpack](https://webpack.js.org/) to build a `bundle.js` that includes the `npm` modules.

Build the worker using:

```bash
# Install dependencies
npm install

# Build with webpack
npm run build
```

The output of the build goes into the `./dist` directory.

## The policy
Using the `example.rego` file build it with:

```bash
opa build -d example.rego 'data.example.allow = true'
```

The input for the policy is going to be a JSON string like:

```json
{
  "fetcher": {},
  "redirect": "manual",
  "headers": {},
  "url": "https://opa4fun.com/",
  "method": "GET",
  "bodyUsed": false,
  "body": null
}
```

The policy is only giving a boolean response for `data.example.allow == true`, in the
`example.rego` policy file that means that it is a `GET` request or a `POST` to `/api`.

The policy itself isn't very exciting, but acts as a proof of concept. With access to
the full `Request` and `Response` objects plus additional Cloudflare API's and headers
you can do some pretty powerful things.


## Uploading the Cloudflare worker

### Prerequisites

Before using workers you must have a Cloudflare account and a site configured to use
Cloudflare. See (https://www.cloudflare.com/)[https://www.cloudflare.com/] for instructions
to get started.

You must also have Workers enabled. See
[https://www.cloudflare.com/products/cloudflare-workers/](https://www.cloudflare.com/products/cloudflare-workers/)

### Test It/Upload Via Editor
Open up the Editor in the Cloudflare dashboard and copy in the `dist/bundle.js` contents.
Then upload the `policy.wasm` file as a "Resource" with the name `REGO_WASM` (needs to match
the variable used in the code).

From the Editor you can configure & deploy the worker.

### Upload via API

Alternatively, and arguably easier (once the worker is stable), they can
be managed via API's.

For the commands shown below either setup some env variables for CF
credentials and account info, or substitute values manually as required. See
https://developers.cloudflare.com/workers/api/ for details. 

```bash
export CF_API_KEY=xxxxxxxxxxxxxxxxxxxxxxxx
export CF_EMAIL=example@gmail.com
export CF_ZONE_ID=yyyyyyyyyyyyyyyyyyyyyyy
```

Below are some helpers to make API calls. The official documentation is the
best source of truth, but these should help get started.

## Create a route

The worker script needs to be associated with a route. Start by creating one,
or checking if some exist already (see below).

> Make sure to substitute the domain name for your application, the example
  below is using a fake url `example.com`

```bash
curl -X POST "https://api.cloudflare.com/client/v4/zones/${CF_ZONE_ID}/workers/filters" \
    -H "X-Auth-Email:${CF_EMAIL}" \
    -H "X-Auth-Key:${CF_API_KEY}" \
    -H "Content-type: application/json" \
    -d '{"pattern": "example.com/*", "enabled": true}'
```

This example just adds a route that will catch pretty much any API call for the site. See
[https://developers.cloudflare.com/workers/api/route-matching/](https://developers.cloudflare.com/workers/api/route-matching/) for more details on configuring the route.

## Get a route

```bash
curl -X GET "https://api.cloudflare.com/client/v4/zones/${CF_ZONE_ID}/workers/filters" \
    -H "X-Auth-Email:${CF_EMAIL}" \
    -H "X-Auth-Key:${CF_API_KEY}"
```
_Initially this should be empty_


## POST a new worker/update the worker

This uses the [Resource Bindings API](https://developers.cloudflare.com/workers/api/resource-bindings/) to upload both a script and wasm binary together, with bindings
for how to load the WebAssembly.

The key being the [metadata_wasm.json](./metadata_wasm.json) file which defines what
variable in the worker context to bind the loaded wasm module to, as well as the
form parts that are the script and assembly files.

```bash
curl -X PUT "https://api.cloudflare.com/client/v4/zones/${CF_ZONE_ID}/workers/script" \
    -H "X-Auth-Email:${CF_EMAIL}" \
    -H "X-Auth-Key:${CF_API_KEY}" \
    -F 'metadata=@metadata_wasm.json;type=application/json' \
    -F 'wasmpolicy=@policy.wasm;type=application/wasm' \
    -F 'script=@./dist/bundle.js;type=application/javascript'
```

## GET the current worker

```bash
curl -X GET "https://api.cloudflare.com/client/v4/zones/${CF_ZONE_ID}/workers/script" \
    -H "X-Auth-Email:${CF_EMAIL}" \
    -H "X-Auth-Key:${CF_API_KEY}"
```

_Initially this should be empty_


# Exercising the policy

For the example policy you can verify it works by making requests like:

```bash
curl -svL -X GET https://${DOMAIN_NAME}/
```
With a successful `200` response (depends on the backing HTTP service).


```bash
curl -svL -X PUT https://${DOMAIN_NAME}/
```

With a `403` response like:

```json
{
  "error": "Not allowed by policy"
}
```