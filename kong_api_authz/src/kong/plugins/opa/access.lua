local cjson_safe = require "cjson.safe"
local http = require "resty.http"
local jwt = require "resty.jwt"

-- string interpolation with named parameters in table
local function interp(s, tab)
    return (s:gsub('($%b{})', function(w) return tab[w:sub(3, -2)] or w end))
end

-- slice a list
local function slice(list, from, to)
    local sliced_results = {};
    for i=from, to do
        table_insert(sliced_results, list[i]);
    end;
    return sliced_results;
end

-- split string based on delimiter
local function split(s, delimiter)
    local result = {};
    for match in (s..delimiter):gmatch("(.-)"..delimiter) do
        table_insert(result, match);
    end
    return result
end

-- query "Get a Document (with Input)" endpoint from the OPA Data API
local function getDocument(input, conf)
    -- serialize the input into a string containing the JSON representation
    local json_body = assert(cjson_safe.encode({ input = input }))

    local opa_uri = interp("${protocol}://${host}:${port}/${base_path}/${decision}", {
        protocol = conf.server.protocol,
        host = conf.server.host,
        port = conf.server.port,
        base_path = conf.policy.base_path,
        decision = conf.policy.decision
    })

    local res, err = http.new():request_uri(opa_uri, {
        method = "POST",
        body = json_body,
        headers = {
          ["Content-Type"] = "application/json",
        },
        keepalive_timeout = conf.server.connection.timeout,
        keepalive_pool = conf.server.connection.pool
    })

    if err then
        error(err) -- failed to request the endpoint
    end

    -- deserialise the response into a Lua table
    return assert(cjson_safe.decode(res.body))
end

-- module
local _M = {}

function _M.execute(conf)
    local authorization = ngx.var.http_authorization

    -- decode JWT token
    local token = {}
    if authorization and string.find(authorization, "Bearer") then
        local encoded_token = authorization:gsub("Bearer ", "")
        token = jwt:load_jwt(encoded_token)
    end

    local list_path = split(ngx.var.upstream_uri, "/")
    local split_path = slice(list_path, 2, #list_path)

    -- input document that will be send to opa
    local input = {
        token = token,
        method = ngx.var.request_method,
        path = ngx.var.upstream_uri,
        split_path = split_path,
    }

    local status, res = pcall(getDocument, input, conf)

    if not status then
        kong.log.err("Failed to get document: ", res)
        return kong.response.exit(500, [[{"message":"Oops, something went wrong"}]])
    end

    -- when the policy fail, 'result' is omitted
    if not res.result then
        kong.log.info("Access forbidden")
        return kong.response.exit(403, [[{"message":"Access Forbidden"}]])
    end

    -- access allowed
    kong.log.debug(interp("Access allowed to ${method} ${path} for user ${subject}", {
        method = input.method,
        path = input.path,
        subject = (token.payload and token.payload.sub or 'anonymous')
    }))
end

return _M