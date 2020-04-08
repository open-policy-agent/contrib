FROM alpine:latest

LABEL org.openpolicyagent.name="openpolicyagent/kong-plugin-opa-devcontainer"
LABEL org.openpolicyagent.description="Development container for kong-plugin-opa"

# Version compatible with Kong 2.0.x
ENV LUA_VERSION=5.1.5
ENV LUAROCKS_VERSION=3.3.1

# install development tools
RUN apk add --no-cache --virtual build-essential \
    make gcc libc-dev readline-dev curl unzip openssl

# build and install Lua
RUN wget -O - http://www.lua.org/ftp/lua-${LUA_VERSION}.tar.gz | tar -zxf - \
    && cd lua-${LUA_VERSION}/ \
    && make linux test \
    && make install

# download and unpack the LuaRocks tarball
RUN wget --no-check-certificate -O - https://luarocks.org/releases/luarocks-${LUAROCKS_VERSION}.tar.gz | tar -zxpf - \
    && cd luarocks-${LUAROCKS_VERSION}/ \
    && ./configure \
    && make build \
    && make install

# install luacheck
RUN luarocks install luacheck

COPY entrypoint.sh /
ENTRYPOINT [ "/entrypoint.sh" ]