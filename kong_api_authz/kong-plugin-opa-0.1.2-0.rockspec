rockspec_format = "3.0"
package = "kong-plugin-opa"
version = "0.1.2-0"
source = {
   url = "git+https://github.com/open-policy-agent/contrib.git",
   tag = "v0.1.2",
}
description = {
   summary = "Integrate the Open Policy Agent (OPA) with Kong API Gateway for API access management",
   detailed = [[
      see https://github.com/open-policy-agent/contrib/tree/master/kong_plugin_opa for more information
   ]],
   homepage = "https://github.com/open-policy-agent/contrib/tree/master/kong_plugin_opa",
   issues_url = "https://github.com/open-policy-agent/contrib/issues",
}
dependencies = {
   "lua-cjson",
   "lua-resty-http",
   "lua-resty-jwt",
}
test_dependencies = {
   "luacov",
   "luacheck",
}
test = {
   type = "busted",
   flags = { "-o", "gtest" },
}
build = {
   type = "builtin",
   modules = {
      ["kong.plugins.opa.access"] = "src/kong/plugins/opa/access.lua",
      ["kong.plugins.opa.handler"] = "src/kong/plugins/opa/handler.lua",
      ["kong.plugins.opa.schema"] = "src/kong/plugins/opa/schema.lua",
   },
}
