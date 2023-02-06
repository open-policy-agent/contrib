-- define kong-plugin-opa configuration
local conf = {
  server = {
    protocol = "https",
    host = "opa",
    port = 8181,
    connection = {
      timeout = 5,
      pool = 1,
    }
  },
  policy = {
    base_path = "v1/data",
    decision = "opa/example/allow",
  }
}

-- mock incoming request
_G.ngx = {var = {}, req = {}}
_G.ngx.var.http_authorization = "Bearer JWT_TOKEN"

-- mock kong loggers and responses on forbidden access or error
_G.kong = {log = {}, response = {}}
_G.kong.response.exit = function(status, body)
  -- value returned by functions on error
  return status
end

local emptyFunction = function() end
_G.kong.log.debug = emptyFunction
_G.kong.log.info = emptyFunction
_G.kong.log.err = emptyFunction
_G.ngx.req.get_headers = emptyFunction

describe("opa:access", function()
  local access

  setup(function()
    -- add mocked modules to path
    package.path = package.path..";spec/?.lua"
    -- override module loader to use fake modules
    package.loaded["resty.http"] = require("__mocks__.resty.http")
    package.loaded["resty.jwt"] = require("__mocks__.resty.jwt")
    -- load opa:access module
    access = require("kong.plugins.opa.access")
  end)

  before_each(function()
    -- reset request method and path
    _G.ngx.var.request_method = "GET"
    _G.ngx.var.upstream_uri = "/api/endpoint"
  end)

  it("allow access", function()
    local res = access.execute(conf)
    assert.is_nil(res)
  end)

  it("returns 403 when request is forbidden", function()
    -- using "/not_allowed" will force fake-http module to mock a negative decision
    _G.ngx.var.upstream_uri = "/api/endpoint/not_allowed"
    local res = access.execute(conf)
    assert.is_true(res == 403)
  end)

  it("returns 500 on OPA server error or when not reachable", function()
    -- using "/error" will force fake-http module to mock an error
    _G.ngx.var.upstream_uri = "/api/endpoint/error"
    local res = access.execute(conf)
    assert.is_true(res == 500)
  end)

  it("sends a request to the server defined in the configuration", function()
    local match = require("luassert.match")
    -- spy on the fake_http module to catch the call made to opa
    local fake_http = require("resty.http")
    spy.on(fake_http, "request_uri")
    local server = conf.server
    local expected_uri = server.protocol .. "://" .. server.host .. ":" .. server.port .. "/" .. conf.policy.base_path
    access.execute(conf)
    assert.spy(fake_http.request_uri).was.called(1) -- was called once
    assert.spy(fake_http.request_uri).was.called_with(match._, match.has_match(expected_uri) , match._)
  end)
end)
