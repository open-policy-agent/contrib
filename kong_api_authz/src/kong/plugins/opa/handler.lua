local access = require "kong.plugins.opa.access"

local OpaHandler = {
  -- This plugin should be executed after any authentication plugins enabled on the Service or Route.
  -- The priority is set to execute this plugin after the response-ratelimiting plugin:
  -- https://docs.konghq.com/2.0.x/plugin-development/custom-logic/#plugins-execution-order
  PRIORITY = 899, -- set the plugin priority, which determines plugin execution order
  VERSION = "0.1"
}

-- Execute OPA query on every request from a client, before it is being proxied to the upstream service
function OpaHandler:access(conf)
  kong.log.debug("executing plugin \"opa\": access")
  access.execute(conf)
end

-- return handler
return OpaHandler
