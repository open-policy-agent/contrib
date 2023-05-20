local _M = {}

local mt = { __index = _M }

function _M.new(_)
    return setmetatable({ sock = "", keepalive = true, timeouts = {} }, mt)
end

function _M.set_timeouts(self, connect_timeout, send_timeout, read_timeout)
    -- store timeouts for later use in simulation of timeout scenarios
    self.timeouts = {
        connect_timeout = connect_timeout,
        send_timeout = send_timeout,
        read_timeout = read_timeout,
    }
end

function _M.request_uri(self, uri, params)
    local res = {}

    -- mock a timeout
    if (self.timeouts.read_timeout == 999) then
        return nil, "error"
    end

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
