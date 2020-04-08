#!/bin/sh

######################################################################
# Entrypoint script to init the lua environment and build the plugin

# initialize the project environment
luarocks init

# build and install the kong-plugin-opa
rockspec=$(ls kong-plugin-opa-*.rockspec)
luarocks make $rockspec

# execute command
$@