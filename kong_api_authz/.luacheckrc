-- Configuration file for LuaCheck
-- see: https://luacheck.readthedocs.io/en/stable/
--
-- To run do: `luacheck .` from root directory

std             = "ngx_lua"
unused_args     = false
redefined       = false

globals = {
    "_KONG",
    "kong",
    "ngx.IS_CLI",
}

not_globals = {
    "string.len",
    "table.getn",
}

exclude_files = {
    "lua_modules",
    ".luarocks"
}

files["spec/**/*.lua"] = {
    std = "ngx_lua+busted",
}
