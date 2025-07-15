local cjson_safe = require "cjson.safe"
local http = require "resty.http"
local jwt = require "resty.jwt"
local helpers = require "kong.plugins.opa.helpers"

-- module
local _M = {}

-- query "Get a Document (with Input)" endpoint from the OPA Data API
local function getDocument(input, conf)
    -- serialize the input into a string containing the JSON representation
    local json_body = assert(cjson_safe.encode({ input = input }))

    local opa_uri = helpers.interp("${protocol}://${host}:${port}/${base_path}/${decision}", {
        protocol = conf.server.protocol,
        host = conf.server.host,
        port = conf.server.port,
        base_path = conf.policy.base_path,
        decision = conf.policy.decision
    })

    local resty_http = http.new()
    resty_http:set_timeouts(conf.server.connection.connect_timeout,
        conf.server.connection.send_timeout,
        conf.server.connection.read_timeout)
    local res, err = resty_http:request_uri(opa_uri, {
        method = "POST",
        body = json_body,
        headers = {
            ["Content-Type"] = "application/json",
        },
        keepalive_timeout = conf.server.connection.timeout,
        keepalive_pool = conf.server.connection.pool
    })

    if err then
        kong.log.err("Error calling OPA for URI " .. opa_uri ..
            " ; error: ", err)
        error(err) -- failed to request the endpoint
    end

    -- deserialise the response into a Lua table
    return assert(cjson_safe.decode(res.body))
end

function _M.execute(conf)
    local authorization = ngx.var.http_authorization

    -- decode JWT token
    local token = {}
    if authorization and string.find(authorization, "Bearer") then
        local encoded_token = authorization:gsub("Bearer ", "")
        token = jwt:load_jwt(encoded_token)
    end

    -- input document that will be send to opa
    local input = {
        token = token,
        method = ngx.var.request_method,
        path = ngx.var.upstream_uri,
        headers = helpers.filterHeaders(ngx.req.get_headers(), conf.document and conf.document.include_headers)
    }

    local status, res = pcall(getDocument, input, conf)

    if not status then
        kong.log.err("Failed to get document: ", res)
        return kong.response.exit(500, [[{"message":"An error occurred while authorizing the request."}]])
    end

    -- Handle both boolean and object formats
    if type(res.result) == "boolean" then
        if not res.result then
            kong.log.info("Access forbidden - boolean result is false")
            return kong.response.exit(403, [[{"message":"Access Forbidden"}]])
        end
    elseif type(res.result) == "table" then
        if not res.result.allow then
            local status_code = res.result.status or 403
            local message = res.result.message or "Access Forbidden"
            kong.log.info("Access forbidden - object result.allow is false/nil")
            return kong.response.exit(status_code, cjson_safe.encode({message = message}))
        end
    else
        kong.log.err("Invalid result type: " .. type(res.result))
        return kong.response.exit(500, [[{"message":"Invalid policy result format"}]])
    end

    -- access allowed
    kong.log.debug(helpers.interp("Access allowed to ${method} ${path} for user ${subject}", {
        method = input.method,
        path = input.path,
        subject = (token.payload and token.payload.sub or 'anonymous')
    }))
end

return _M
