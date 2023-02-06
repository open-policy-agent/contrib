-- module for helper functions
-- separation into a separate module de-clutters the main access.lua code and enables
-- simplified unit testing


-- module
local _M = {}

function _M.filterHeaders(headers, wanted_headers)
    -- per https://openresty-reference.readthedocs.io/en/latest/Lua_Nginx_API/
    -- "Since the 0.6.9 release, all the header names in the Lua table returned
    -- are converted to the pure lower-case form by default, unless the raw
    -- argument is set to true (default to false)."
    -- So we need to convert the requested header names to lower case too.
    local filtered_headers = {}
    if wanted_headers and headers then
        for _, wanted_header in ipairs(wanted_headers) do
            local lower_key = string.lower(wanted_header)
            local value = headers[lower_key]
            if value then
                filtered_headers[lower_key] = value
            end
        end
    end
    return filtered_headers
end

-- string interpolation with named parameters in table
function _M.interp(s, tab)
    return (s:gsub('($%b{})', function(w) return tab[w:sub(3, -2)] or w end))
end

return _M
