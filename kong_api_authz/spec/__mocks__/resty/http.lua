local _M = {}

local mt = { __index = _M }

function _M.new(_)
    return setmetatable({ sock = "", keepalive = true }, mt)
end

function _M.request_uri(self, uri, params)
    local res = {}

    -- mock 400/500 OPA responses
    if (params.body:find("/error")) then
        return nil, "error"
    end

    -- mock 200 OPA responses
    if (params.body:find("/not_allowed")) then
        -- when the document does not evaluate to true,
        -- the response will not contain the result property
        res.body = '{}'
    else
        -- when the document evaluates to true,
        -- an element is produced in the result
        res.body = '{ "result": "true" }'
    end

    return res, nil
end

return _M